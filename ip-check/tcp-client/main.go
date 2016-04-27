package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func main() {
	var addr = flag.String("a", "stacktitan.com:5555", "address to connect to")
	flag.Parse()
	conn, err := net.Dial("tcp", *addr)
	if err != nil {
		log.Fatal(err)
	}
	buf := make([]byte, 1024)
	for {
		_, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
			break
		}
	}
	fmt.Printf("%s\n", strings.Split(string(buf), ":")[0])
}
