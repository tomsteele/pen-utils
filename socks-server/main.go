package main

import (
	"flag"
	"log"

	"os"

	"github.com/things-go/go-socks5"
)

func main() {

	addr := flag.String("addr", ":10800", "addr to listen on")
	flag.Parse()

	server := socks5.NewServer(
		socks5.WithLogger(socks5.NewLogger(log.New(os.Stdout, "socks5: ", log.LstdFlags))),
	)
	if err := server.ListenAndServe("tcp", *addr); err != nil {
		panic(err)
	}

}
