/* open-things opens a list of files or urls */
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/skratchdot/open-golang/open"
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
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		open.Run(scanner.Text())
	}
}
