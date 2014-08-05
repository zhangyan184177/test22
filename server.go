package main

import (
	"./memcache"
	"log"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", "0.0.0.0:6666")
	if err != nil {
		panic("listening err:"+err.Error())	
	}
	data := make(map[string]string)
	expir := make(map[string]int64)

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic("accept err:"+err.Error())	
		}
		log.Printf("accept connection: %s", conn.RemoteAddr())

		go memcache.Handler(conn, data, expir)
	}
}
