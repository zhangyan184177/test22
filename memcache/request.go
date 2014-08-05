package memcache

import (
	"strings"
	"errors"
	"bufio"
	"io"
	"time"
	"strconv"
)
type Request struct {
	cmd string
	keys string
	value string
	length int
	result string
	delay int64 
	interval int64 
}

func ReadData(reader io.Reader, data map[string]string, expir map[string]int64) (*Request, error) {
	req := new(Request)
	buf := bufio.NewReader(reader)
	line, _, err := buf.ReadLine()
	if err != nil && len(line) == 0 {
		return req, io.EOF	
	}
	params := strings.Fields(string(line))
	cmd := params[0]
	switch cmd {
		case SET, ADD, REPLACE:
			req.cmd = cmd
			req.keys = params[1]
			exist := ExistMap(params[1], data)
			if cmd == ADD && exist == true || cmd == REPLACE && exist == false {
				req.result = NotStored
				return req, errors.New("add a exised key or replace a not exised key")
			}
			now := time.Now()
			req.interval, _ = strconv.ParseInt(params[3], 10, 64)
			expir_time := now.Add(time.Duration(req.interval) * time.Second)
			tmp_time := time.Date(expir_time.Year(), expir_time.Month(), expir_time.Day(), expir_time.Hour(), expir_time.Minute(), expir_time.Second(), 0, time.Local)
			save_time := tmp_time.Unix()
			if req.interval != 0 && save_time > expir[req.keys] {
				expir[req.keys] = save_time
				go DataExpir(req, data, expir)
			}
			req.value = params[2]
			data[req.keys] = req.value
			
			req.result = Stored
		case GET:
			req.cmd = cmd
			req.keys = params[1]
			exist := ExistMap(params[1], data)
			if exist == false {
				req.result = NotFound
				return req, errors.New("get a not exised key")
			}
			req.value = data[req.keys]	
		case DELETE:
			req.cmd = cmd
			req.keys = params[1]
			exist := ExistMap(params[1], data)
			if exist == false {
				req.result = NotFound
				return req, errors.New("delete a not exised key")
			}
			if len(params) == 3 {
				req.delay, _ = strconv.ParseInt(params[2], 10, 64)
				go DelayDelete(req, data)
			} else {
				delete(data, req.keys)	
			}
			req.result = Deleted
		case FLUSH_ALL:
			req.cmd = cmd
			for i, _ := range data {
				delete(data, i)	
			}
			req.result = OK
		default:
			req.cmd = INVAILD 
			return req, errors.New("cmd is invaild")
	}
	return req, err
}

func ExistMap(key string, data map[string]string) bool {
	is_exist := false
	for i, _ := range data {
		if i == key {
			is_exist = true
			break
		}	
	}
	return is_exist
}

func DelayDelete(req *Request, data map[string]string) {
	delay_timer := time.NewTimer(time.Duration(req.delay) * time.Second)
	<-delay_timer.C
	exist := ExistMap(req.keys, data)
	if exist == true {
		delete(data, req.keys)	
	}
	return
}

func DataExpir(req *Request, data map[string]string, expir map[string]int64) {
	expir_timer := time.NewTimer(time.Duration(req.interval) * time.Second)
	<-expir_timer.C
	now := time.Now().Unix()
	if now >= expir[req.keys] {
		exist := ExistMap(req.keys, data)
		if exist == true {
			delete(data, req.keys)	
		}
	}
	return
}
