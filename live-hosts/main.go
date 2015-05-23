/*live-hosts reads an nmap xml file from stdin or as an argument and shows the hosts that are alive*/
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
			fmt.Println(h.Address[0].Addr)
		}
	}
}
