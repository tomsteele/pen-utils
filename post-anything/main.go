/* post-anything allows you to send HTTP POST/GET to anything, anywhere.
supports CORS, and can be used as a static file server*/
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"

	"github.com/mholt/binding"
)

type fileRequest struct {
	File *multipart.FileHeader
}

func (f *fileRequest) FieldMap(req *http.Request) binding.FieldMap {
	return binding.FieldMap{
		&f.File: "file",
	}
}

func main() {
	listen := flag.String("listen", "", "interface to listen on e.g. 0.0.0.0:8000")
	dir := flag.String("dir", "", "serve static files from here, will use url of /static/*")
	flag.Parse()
	if *listen == "" {
		log.Fatal("listen is required")
	}

	http.HandleFunc("/sendfile", func(w http.ResponseWriter, req *http.Request) {
		fReq := &fileRequest{}
		binding.MaxMemory = 104857600000
		if errs := binding.Bind(req, fReq); errs.Len() != 0 {
			w.WriteHeader(500)
			return
		}
		if fReq.File == nil {
			w.WriteHeader(400)
			return
		}
		fh, err := fReq.File.Open()
		if err != nil {
			w.WriteHeader(500)
			return
		}
		defer fh.Close()
		fh2, err := os.Create(fReq.File.Filename)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		defer fh2.Close()

		if _, err := io.Copy(fh2, fh); err != nil {
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(200)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		switch req.Method {
		case "GET":
			logMap(req.Form)
		case "POST":
			logMap(req.Form)
		case "OPTIONS":
			w.Header().Add("Access-Control-Allow-Origin", "*")
		}
		w.WriteHeader(200)
	})

	if *dir != "" {
		http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(*dir))))
	}

	log.Fatal(http.ListenAndServe(*listen, nil))
}

func logMap(m url.Values) {
	fmt.Println("Request Values")
	for k, v := range m {
		fmt.Println(k, v)
	}
}
