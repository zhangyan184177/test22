package memcache

import (
	"strconv"
	"strings"
	"errors"
	"bufio"
	"time"
	"io"
)

type Request struct {
	cmd string
	key string
	value []byte
	length int
	result string
	delay int64
	interval int64
}

func (req *Request) ReadData(read io.Reader,
			data map[string][]byte, expir map[string]int64) error {
	reader := bufio.NewReaderSize(read, NormalReaderSize)
	line, err := reader.ReadBytes('\n')
	if err != nil {
		return err
	}
	if len(line) < 2 || line[len(line)-1] != '\n' || line[len(line)-2] != '\r' {
		req.result = Invaild
		return errors.New("command format is invaild")
	}
	
	params := strings.Fields(string(line))
	if len(params) < 1 {
		req.result = Invaild
		return errors.New("lack of instruction")
	}
	req.cmd = params[0]
	switch req.cmd {
		case SET, ADD, REPLACE:
			if len(params) != StoreLen {
				req.result = Invaild
				return errors.New("commond format of store is wrong")
			}
			req.key = params[1]
			_, exist := data[req.key]
			if req.cmd == ADD && exist == true || req.cmd == REPLACE && exist == false {
				req.result = NotStored
				return errors.New("add a exised key or replace a not exised key")
			}
			now := time.Now()
			req.interval, _ = strconv.ParseInt(params[3], 10, 64)
			expir_time := now.Add(time.Duration(req.interval) * time.Second)
			tmp_time := time.Date(expir_time.Year(), expir_time.Month(), expir_time.Day(),
				expir_time.Hour(), expir_time.Minute(), expir_time.Second(), 0, time.Local)
			save_time := tmp_time.Unix()
			if req.interval != 0 && save_time > expir[req.key] {
				expir[req.key] = save_time
				go DataExpir(req, data, expir)
			}
			req.value = []byte(params[2])
//			data[req.key] = req.value
			data[req.key] = make([]byte, len(req.value))
			copy(data[req.key], req.value)
			req.result = Stored
		case GET:
			if len(params) != GetLen {
				req.result = Invaild
				return errors.New("commond format of get is wrong")
			}
			req.key = params[1]
			_, exist := data[req.key]
			if exist == false {
				req.result = NotFound
				return errors.New("get a not exised key")
			}
//			req.value = data[req.key]
			req.value = make([]byte, len(data[req.key]))
			copy(req.value, data[req.key])
		case DELETE:
			if len(params) != DelNowLen && len(params) != DelDelayLen {
				req.result = Invaild
				return errors.New("commond format of delete is wrong")
			}
			req.key = params[1]
			_, exist := data[req.key]
			if exist == false {
				req.result = NotFound
				return errors.New("delete a not exised key")
			}
			if len(params) == DelDelayLen {
				req.delay, _ = strconv.ParseInt(params[2], 10, 64)
				go DelayDelete(req, data)
			} else {
				delete(data, req.key)
			}
			req.result = Deleted
		case FLUSH_ALL:
			if len(params) != FlushAllLen {
				req.result = Invaild
				return errors.New("commond format of flush_all is wrong")
			}
			for i, _ := range data {
				delete(data, i)
			}
			req.result = OK
		default:
			return errors.New("cmd is invaild")
	}
	return err
}

func DelayDelete(req *Request, data map[string][]byte) {
	delay_timer := time.NewTimer(time.Duration(req.delay) * time.Second)
	<-delay_timer.C
	delete(data, req.key)
	return
}

func DataExpir(req *Request, data map[string][]byte, expir map[string]int64) {
	expir_timer := time.NewTimer(time.Duration(req.interval) * time.Second)
	<-expir_timer.C
	now := time.Now().Unix()
	if now >= expir[req.key] {
		delete(data, req.key)
	}
	return
}
