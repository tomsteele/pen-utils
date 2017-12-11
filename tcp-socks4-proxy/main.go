package main

import (
	"flag"
	"io"
	"log"
	"net"

	"github.com/eahydra/socks"
)

func checkAndFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	socksAddr := flag.String("socks-addr", "", "the host:port for your socks proxy")
	dstAddr := flag.String("dest-addr", "", "the host:port for your destination")
	listenAddr := flag.String("listen-addr", "", "the host:port for listen on")
	flag.Parse()

	client, err := socks.NewSocks4Client("tcp", *socksAddr, "", socks.Direct)
	checkAndFatal(err)

	listen, err := net.Listen("tcp", *listenAddr)
	checkAndFatal(err)

	for {
		y, err := listen.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go func(z net.Conn) {
			x, err := client.Dial("tcp", *dstAddr)
			checkAndFatal(err)
			go io.Copy(x, z)
			io.Copy(z, x)
			x.Close()
			y.Close()
		}(y)
	}
}
