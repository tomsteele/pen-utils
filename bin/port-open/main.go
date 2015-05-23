/*port-open reads an nmap xml file from stdin or as an argument and shows the hosts and port with that port open*/
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/lair-framework/go-nmap"
)

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
func main() {
	var file *os.File
	var err error
	format := flag.String("f", "%s:%s", "ip and port format string")
	port := flag.String("p", "80", "port to be open")
	flag.Parse()
	if len(flag.Args()) >= 1 {
		file, err = os.Open(flag.Arg(0))
		checkError(err)
	} else {
		file = os.Stdin
	}
	data, err := ioutil.ReadAll(file)
	checkError(err)
	n, err := nmap.Parse(data)
	checkError(err)

	for _, h := range n.Hosts {
		if h.Status.State == "up" {
			ip := h.Address[0].Addr
			for _, p := range h.Ports {
				if p.PortId == *port && p.State.State == "open" {
					fmt.Printf(*format, ip, p.PortId)
					fmt.Println()
				}
			}
		}
	}
}
