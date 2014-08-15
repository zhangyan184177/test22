package memcache

import (
	"net"
	"log"
	"time"
)

func Handler(conn net.Conn, datachan chan *Request) {
	req := Request{}
	err_read := req.ReadData(conn)
	if err_read != nil {
		log.Println("read data from client failed:", err_read)
	}
	operatemap(&req, datachan)

	rsp := Response{}
	rsp.cmd = req.cmd
	rsp.key = req.key
	rsp.value = req.value
	rsp.result = req.result
	err_write := rsp.WriteData(conn)
	if err_write != nil {
		log.Println("write data to client failed:", err_write)
	}
	return
}

func operatemap(req *Request, datachan chan *Request) {
	datachan <- req
	<-req.clientchan
}

func SyncData(datachan chan *Request, data map[string][]byte,
									expir map[string]int64) {
	for {
		req := <-datachan
		switch req.cmd {
			case SET, ADD, REPLACE:
				_, exist := data[req.key]
				if req.cmd == ADD && exist == true ||
						req.cmd == REPLACE && exist == false {
					req.result = NotStored
					log.Println("add a exised key"+
							" or replace a not exised key")
				} else {
					if req.interval != 0 {
						updateexpir(expir, req.interval, req.key)
					}
					data[req.key] = req.value
					req.result = Stored
				}
			case GET:
				now := time.Now().Unix()
				_, exist := data[req.key]
				if exist == false || expir[req.key] != 0 && now >= expir[req.key] {
					req.result = NotFound
					log.Println("get a not exised key")
				} else {
					req.value = data[req.key]
				}
			case DELETE:
				_, exist := data[req.key]
				if exist == false {
					req.result = NotFound
					log.Println("delete a not exised key")
				} else {
					if req.delay != 0 {
						updateexpir(expir, req.delay, req.key)
					}
					delete(data, req.key)
					req.result = Deleted
				}
			case FLUSH_ALL:
				for i, _ := range data {
					delete(data, i)
				}
				req.result = OK
		}
		req.clientchan <- true
	}
}

func updateexpir(expir map[string]int64, intertime int64, key string) {
	now := time.Now()
	expir_time := now.Add(time.Duration(intertime) * time.Second)
	expir_time = time.Date(expir_time.Year(), expir_time.Month(),
			expir_time.Day(), expir_time.Hour(), expir_time.Minute(),
					expir_time.Second(), 0, time.Local)
	if expir_time.Unix() > expir[key] {
		expir[key] = expir_time.Unix()
	}
}
