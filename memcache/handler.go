package memcache

import (
	"net"
	"log"
	"time"
)
type Request struct {
	cmd string
	key string
	value []byte
	length int
	result string
	delay int64
	interval int64
	delayflag bool
	intervalflag bool
	clientchan chan bool
}

func Handler(conn net.Conn, datachan chan *Request,
								expirchan chan *Request) {
	req := Request{}
	req.clientchan = make(chan bool)
	err_read := req.ReadData(conn, datachan, expirchan)
	if err_read != nil {
		log.Println("read data from client failed:", err_read)
	}

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

func SyncData(datachan chan *Request, data map[string][]byte) {
	log.Println("sync")
	for {
		req := <-datachan
		switch req.cmd {
			case SET, ADD, REPLACE:
				if req.intervalflag {
					delete(data, req.key)
					break
				} else {
					_, exist := data[req.key]
					if req.cmd == ADD && exist == true ||
							req.cmd == REPLACE && exist == false {
						req.result = NotStored
						log.Println("add a exised key"+
								" or replace a not exised key")
					} else {
						data[req.key] = req.value
						req.result = Stored
					}
					req.clientchan <- true
				}
			case GET:
				_, exist := data[req.key]
				if exist == false {
					req.result = NotFound
					log.Println("get a not exised key")
				} else {
					req.value = data[req.key]
				}
				req.clientchan <- true
			case DELETE:
				if req.delayflag {
					delete(data, req.key)
					break
				} else {
					_, exist := data[req.key]
					if exist == false {
						req.result = NotFound
						log.Println("delete a not exised key")
					} else {
						delete(data, req.key)
						req.result = Deleted
					}
					req.clientchan <- true
				}
			case FLUSH_ALL:
				for i, _ := range data {
					delete(data, i)
				}
				req.result = OK
				req.clientchan <- true
		}
	}
}

func ExpirData(datachan chan *Request,
		expirchan chan *Request, expir map[string]int64 ) {
	for {
		req := <-expirchan
		now := time.Now()
		expir_time := now.Add(time.Duration(req.interval) * time.Second)
		tmp_time := time.Date(expir_time.Year(), expir_time.Month(),
				expir_time.Day(), expir_time.Hour(), expir_time.Minute(),
								expir_time.Second(), 0, time.Local)
		save_time := tmp_time.Unix()
		if req.interval != 0 && save_time > expir[req.key] {
				expir[req.key] = save_time
				go func() {
					expir_timer :=
						time.NewTimer(time.Duration(req.interval) * time.Second)
					<-expir_timer.C
					now := time.Now().Unix()
					if now >= expir[req.key] {
						req.intervalflag = true
						datachan <- req
					}
				}()
		}
	}
}
