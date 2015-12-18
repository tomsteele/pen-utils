/* post-anything allows you to send HTTP POST/GET to anything, anywhere.
supports CORS, and can be used as a static file server*/
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

func main() {
	listen := flag.String("listen", "", "interface to listen on e.g. 0.0.0.0:8000")
	dir := flag.String("dir", "", "serve static files from here, will use url of /static/*")
	flag.Parse()
	if *listen == "" {
		log.Fatal("listen is required")
	}

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
