/*server-banner takes a url and returns the server header.*/
package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"fmt"
	"net"
	"os"
	"regexp"
	"sort"
)

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var ipreg *regexp.Regexp

func init() {
	ipreg, _ = regexp.Compile(`(?:[0-9]{1,3}\.){3}[0-9]{1,3}`)
}

type ipstr []string

func (i ipstr) Len() int {
	return len(i)
}

func (i ipstr) Swap(j, k int) {
	i[j], i[k] = i[k], i[j]
}

func (i ipstr) Less(j, k int) bool {
	first := net.ParseIP(ipreg.FindString(i[j])).To4()
	second := net.ParseIP(ipreg.FindString(i[k])).To4()
	if first == nil {
		return true
	}
	if second == nil {
		return false
	}
	return binary.BigEndian.Uint32(first) < binary.BigEndian.Uint32(second)
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
	defer file.Close()

	lines := ipstr{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	checkError(scanner.Err())
	sort.Sort(lines)
	for _, i := range lines {
		fmt.Println(i)
	}
}
