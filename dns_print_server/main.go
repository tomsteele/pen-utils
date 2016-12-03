package main

import (
	"fmt"
	"log"

	"github.com/miekg/dns"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	dns.HandleFunc(".", func(w dns.ResponseWriter, req *dns.Msg) {
		for _, q := range req.Question {
			fmt.Println(q.Name)
		}
		dns.HandleFailed(w, req)
	})
	log.Fatal(dns.ListenAndServe(":53", "udp", nil))
}
