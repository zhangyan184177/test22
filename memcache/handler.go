package memcache

import (
	"net"
	"log"
)

func Handler(conn net.Conn, data map[string][]byte, expir map[string]int64) {
	req := Request{}
	err_read := req.ReadData(conn, data, expir)
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
