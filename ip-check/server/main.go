package main

import (
	"log"
	"net"
)

func main() {
	l, err := net.Listen("tcp", ":5555")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
		}
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	ip := conn.RemoteAddr().String()
	log.Printf("Connection from %s\n", ip)
	conn.Write([]byte(ip))
	conn.Close()
}
