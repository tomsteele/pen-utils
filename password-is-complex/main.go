package main

import (
	"os"
	"flag"
	"fmt"
	"bufio"
	"regexp"
)

var matches = []*regexp.Regexp{
	regexp.MustCompile(`[a-z]`),
	regexp.MustCompile(`[A-Z]`),
	regexp.MustCompile(`[0-9]`),
	regexp.MustCompile(`[^0-9a-zA-Z]`),
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func isComplex(p string) bool {
	var typ int
	for _, m := range matches {
		if m.MatchString(p) {
			typ++
		}
	}
	return typ > 2 && len(p) >= 8
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