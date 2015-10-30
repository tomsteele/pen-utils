package main

import (
	"flag"
	"io"
	"log"
	"net"
	"sync"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	srcAddr := flag.String("s", "localhost:3000", "the local host and port to listen on")
	dstAddr := flag.String("d", "", "the destination host and port to listen on")
	flag.Parse()
	if *dstAddr == "" {
		log.Fatal("Missing -d")
	}
	c, err := net.Listen("tcp", *srcAddr)
	check(err)
	for {
		conn, err := c.Accept()
		check(err)
		go func(src net.Conn) {
			w := &sync.WaitGroup{}
			defer src.Close()
			dst, err := net.Dial("tcp", *dstAddr)
			check(err)
			defer dst.Close()
			w.Add(2)
			go func(src, dst net.Conn) {
				io.Copy(src, dst)
				w.Done()
			}(src, dst)
			go func(dst, src net.Conn) {
				io.Copy(dst, src)
				w.Done()
			}(dst, src)
			w.Wait()
		}(conn)
	}
}
