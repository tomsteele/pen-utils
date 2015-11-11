package main

import (
	"crypto/aes"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func aesDecrypt(key, data []byte) ([]byte, error) {
	plain := make([]byte, len(data))

	a, err := aes.NewCipher(key)
	if err != nil {
		return plain, err
	}

	blksize := a.BlockSize()
	if len(data)%blksize != 0 {
		return plain, errors.New("The length of the provided ciphertext is invalid")
	}

	for len(data) > 0 {
		blk := make([]byte, blksize)
		a.Decrypt(blk, data[:blksize])
		data = data[blksize:]
		plain = append(plain, blk...)
	}
	return plain, nil
}

func decode(data string) []byte {
	if x, err := hex.DecodeString(data); err == nil && len(x) > 1 {
		return x
	}
	if x, err := base64.StdEncoding.DecodeString(data); err == nil && len(x) > 1 {
		return x
	}
	return []byte(data)
}

var usage = `
kek, key, and ciphertext must be hex or base64 encoded. File content is assumed to be raw data.

Usage of %s:
	%s <kek> <key> <file>
	%s <kek> <key> <ciphertext>
	%s <key> <file>
	%s <key> <ciphertext>

`

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage, os.Args[0], os.Args[0], os.Args[0], os.Args[0], os.Args[0])
	}
	flag.Parse()
	if len(flag.Args()) < 2 {
		flag.Usage()
		os.Exit(1)
	}
	var kek, key, fileOrCipher string
	switch len(flag.Args()) {
	case 2:
		key = flag.Arg(0)
		fileOrCipher = flag.Arg(1)
	case 3:
		kek = flag.Arg(0)
		key = flag.Arg(1)
		fileOrCipher = flag.Arg(2)
	default:
		flag.Usage()
		os.Exit(1)
	}

	decodedKey := decode(key)
	if kek != "" {
		plainKey, err := aesDecrypt(decode(kek), decodedKey)
		if err != nil {
			log.Fatalf("Error decrypting key with KEK: %s\n", err.Error())
		}
		decodedKey = plainKey
	}

	var cipher []byte
	if _, err := os.Stat(fileOrCipher); err != nil {
		cipher = decode(fileOrCipher)
	} else {
		ciphertext, err := ioutil.ReadFile(fileOrCipher)
		if err != nil {
			log.Fatalf("Error reading file %s: %s\n", fileOrCipher, err.Error())
		}
		cipher = ciphertext
	}

	plain, err := aesDecrypt(decodedKey, cipher)
	if err != nil {
		log.Fatalf("Error decrypting ciphertext: %s\n", err.Error())
	}

	fmt.Println(string(plain))
}
