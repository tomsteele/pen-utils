/*port-open reads an nmap xml file from stdin or as an argument and shows the hosts and port with that port open*/
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"

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

	format := flag.String("f", "%s:%d", "ip and port format string")
	port := flag.String("p", "80", "comma delimited port to search for, can be '-1' to match all ports")
	service := flag.String("s", ".*", "service regex to search for ex: http")
	product := flag.String("pr", ".*", "product regex to search for ex: Apache")
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

	ports := strings.Split(*port, ",")

	for _, h := range n.Hosts {
		if h.Status.State == "up" {
			ip := h.Addresses[0].Addr
			for _, p := range h.Ports {
				sm, serr := regexp.MatchString(*service, p.Service.Name)
				pm, perr := regexp.MatchString(*product, p.Service.Product)
				if (sm && serr == nil) && (pm && perr == nil) && (*port == "-1" || stringInSlice(strconv.Itoa(p.PortId), ports)) && p.State.State == "open" {
					fmt.Printf(*format, ip, strconv.Itoa(p.PortId))
					fmt.Println()
				}
			}
		}
	}
}
