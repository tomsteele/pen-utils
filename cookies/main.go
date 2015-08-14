/*cookie-flags takes a url and returns the cookie set.*/
package main

import (
	"bufio"
	"crypto/tls"
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

func doReq(location string, timeout int) []string {
	cookies := []string{}
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
		return cookies
	}
	for _, c := range res.Cookies() {
		cookies = append(cookies, c.Raw)
	}
	return cookies
}

func main() {
	var file *os.File
	var err error
	timeout := flag.Int("timeout", 1000, "timeout for requests")
	input := flag.String("input", "", "file containing urls")
	flag.Parse()

	if *input != "" {
		file, err = os.Open(*input)
		checkError(err)
	} else if len(flag.Args()) >= 1 {
		location := flag.Arg(0)
		cookies := doReq(location, *timeout)
		for _, c := range cookies {
			fmt.Printf("%s: %s\n", location, c)
		}
		os.Exit(0)
	} else {
		file = os.Stdin
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		location := scanner.Text()
		cookies := doReq(location, *timeout)
		for _, c := range cookies {
			fmt.Printf("%s: %s\n", location, c)
		}
	}
}
