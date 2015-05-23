/*server-banner takes a url and returns the server header.*/
package main

import (
	"bufio"
	"crypto/tls"
	"encoding/csv"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"
)

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func doReq(location string, timeout int) string {
	req, err := http.NewRequest("GET", location, nil)
	checkError(err)
	tr := &http.Transport{
		Dial: func(network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, time.Duration(timeout)*time.Millisecond)
		},
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	res, err := tr.RoundTrip(req)
	if err != nil {
		return ""
	}
	banner := res.Header["Server"]
	if len(banner) >= 1 {
		return banner[0]
	}
	return ""
}

func main() {
	var file *os.File
	var err error
	timeout := flag.Int("timeout", 1000, "timeout for requests")
	input := flag.String("input", "", "file containing urls")
	flag.Parse()
	c := csv.NewWriter(os.Stdout)

	if *input != "" {
		file, err = os.Open(*input)
		checkError(err)
	} else if len(flag.Args()) >= 1 {
		location := flag.Arg(0)
		banner := doReq(location, *timeout)
		if banner == "" {
			os.Exit(1)
		}
		checkError(c.Write([]string{location, banner}))
		c.Flush()
	} else {
		file = os.Stdin
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		location := scanner.Text()
		banner := doReq(location, *timeout)
		if banner == "" {
			continue
		}
		checkError(c.Write([]string{location, banner}))
		c.Flush()
	}
}
