/*count-open reads an nmap xml file from stdin or as an argument and shows the amount of ports open per host*/
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

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func main() {
	var file *os.File
	var err error

	l := flag.Int("l", 1, "Minimum amount of ports open for host to be shown")
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
			ip := h.Addresses[0].Addr
			portsLength := 0
			for _, p := range h.Ports {
				if p.State.State == "open" {
					portsLength++
				}
			}
			if portsLength > *l {
				fmt.Printf("%s:%d\n", ip, portsLength)
			}
		}
	}
}
