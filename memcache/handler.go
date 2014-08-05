package memcache

import (
	"net"
//	"log"
)

func Handler(conn net.Conn, data map[string]string, expir map[string]int64) error {
	req, _ := ReadData(conn, data, expir)
	rsp := new(Response)
	rsp.cmd = req.cmd
	rsp.keys = req.keys
	rsp.value = req.value
	rsp.result = req.result
//	log.Println("req: ", req)
//	log.Println("rsp: ", rsp)
//	log.Println("data: ", data)
	err := rsp.WriteData(conn)
	if err != nil {
		return err	
	}
	return err	
} 
