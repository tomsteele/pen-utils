package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

// explodePorts takes a portstring similar to nmap and returns an array of ports.
func explodePorts(portstring string) ([]int, error) {
	errmsg := "Invalid Port specification"
	ports := []int{}
	portTokens := strings.Split(portstring, ",")
	for _, token := range portTokens {
		switch {
		case strings.Contains(token, "-"):
			sp := strings.Split(token, "-")
			if len(sp) != 2 {
				return ports, errors.New(errmsg)
			}
			start, err := strconv.Atoi(sp[0])
			if err != nil {
				return ports, errors.New(errmsg)
			}
			end, err := strconv.Atoi(sp[1])
			if err != nil {
				return ports, errors.New(errmsg)
			}
			if start > end || start < 1 || end > 65535 {
				return ports, errors.New(errmsg)
			}
			for ; start <= end; start++ {
				ports = append(ports, start)
			}
		case strings.Contains(token, ","):
			sp := strings.Split(token, ",")
			for _, p := range sp {
				i, err := strconv.Atoi(p)
				if err != nil {
					return ports, errors.New(errmsg)
				}
				if i < 1 || i > 65535 {
					return ports, errors.New(errmsg)
				}
				ports = append(ports, i)
			}
		default:
			i, err := strconv.Atoi(token)
			if err != nil {
				return ports, errors.New(errmsg)
			}
			if i < 1 || i > 65535 {
				return ports, errors.New(errmsg)
			}
			ports = append(ports, i)
		}
	}
	return ports, nil
}

func main() {

	if len(os.Args) < 3 {
		fmt.Printf("%s <ip> <portstring>\n", os.Args[0])
		os.Exit(1)
	}

	host := os.Args[1]
	portStr := os.Args[2]

	ports, err := explodePorts(portStr)
	if err != nil {
		fmt.Printf("Error parsing port string: Error: %s", err.Error())
		os.Exit(1)
	}

	wg := sync.WaitGroup{}
	wg.Add(len(ports))

	// We'll need a channel for passing the port to be scanned to our pool of goroutines.
	port := make(chan int)
	fmt.Printf("Host %s\n", host)

	// Here we create 50 goroutines.
	for l := 0; l < 50; l++ {
		go func() {
			// port is a channel, by calling range this goroutine will wait for a port
			// to be sent over the channel.
			for p := range port {
				if conn, err := net.Dial("tcp", host+":"+strconv.Itoa(p)); err == nil {
					fmt.Printf("%d open\n", p)
					conn.Close()
				}
				wg.Done()
			}
		}()
	}

	for _, p := range ports {
		port <- p
	}

	wg.Wait()
	close(port)
	fmt.Println("done")
}
