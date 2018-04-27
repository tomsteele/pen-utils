package main

import (
	"os"
	"flag"
	"fmt"
	"bufio"
	"regexp"
)

var (
	lower  = regexp.MustCompile(`[a-z]`)
	upper  = regexp.MustCompile(`[A-Z]`)
	number = regexp.MustCompile(`[0-9]`)
	marks  = regexp.MustCompile(`[^0-9a-zA-Z]`)
)

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func isComplex(p string) bool {

	if len(p) < 8 {
		return false
	}

	var typ int

	if lower.MatchString(p) {
		typ++
	}
	if upper.MatchString(p) {
		typ++
	}
	if number.MatchString(p) {
		typ++
	}
	if marks.MatchString(p) {
		typ++
	}
	return typ > 2
}


func main() {
	var file *os.File
	var err error
	input := flag.String("input", "", "file containing passwords")
	flag.Parse()

	if *input != "" {
		file, err = os.Open(*input)
		checkError(err)
	} else {
		file = os.Stdin
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		p := scanner.Text()
		if isComplex(p) {
			fmt.Println(p)
		}
	}
}