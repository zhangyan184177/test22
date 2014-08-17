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

	data := make(map[string][]byte)
	expir := make(map[string]int64)
	
	datachan := make(chan *memcache.Request)
	memcache.InitMapUpdateRequest();
	
	for {
		conn, err := ln.Accept()
		if err != nil {
			panic("accept err:"+err.Error())
		}
		log.Printf("accept connection: %s", conn.RemoteAddr())
	
		go memcache.Handler(conn, datachan)
	}
}
