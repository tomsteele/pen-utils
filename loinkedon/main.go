package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"encoding/base64"
	"fmt"
	"strings"
)

type L struct {
	Included []struct {
		FirstName string `json:"firstName,omitempty"`
		LastName string `json:"lastName,omitempty"`
		Occupation string `json:"occupation,omitempty"`
	} `json:"included"`
}

type RequestResponse struct {
	Response struct {
		Headers struct {
			ContentType string `json:"Content-Type"`
		}`json:"headers"`
		Path string `json:"path"`
		Body string `json:"body"`
	} `json:"response"`
}

func main() {
	addr := flag.String("addr", "localhost:3001", "address to listen on")
	flag.Parse()

	http.HandleFunc("/response", func(w http.ResponseWriter, r *http.Request) {
		defer w.Write([]byte("ok"))
		var req RequestResponse
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("Error: %s\n", err.Error())
			return
		}
		if !strings.Contains(req.Response.Headers.ContentType, "application/vnd.linkedin.normalized+json") {
			return
		}
		buff, err := base64.StdEncoding.DecodeString(req.Response.Body)
		if err != nil {
			log.Printf("Error: %s\n", err.Error())
			return
		}
		l := L{}
		if err := json.Unmarshal(buff, &l); err != nil {
			log.Printf("Error: %s\n", err.Error())
			return
		}
		for _, i := range l.Included {
			if i.FirstName == "" {
				continue
			}
			fmt.Printf("first_name: %s last_name: %s title: %s\n", i.FirstName, i.LastName, i.Occupation)
		}

	})

	log.Fatal(http.ListenAndServe(*addr, nil))
}
