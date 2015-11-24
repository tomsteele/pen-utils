/* bing-pastebin searches bing for pastebins associated with a domain */
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type bingMessage struct {
	D bingResults `json:"D"`
}

type bingResults struct {
	Results []bingResult `json:"Results"`
}

type bingResult struct {
	Metadata    bingMetadata `json:"__Metadata"`
	ID          string       `json:"id"`
	Title       string       `json:"Title"`
	Description string       `json:"Description"`
	DisplayURL  string       `json:"DisplayUrl"`
	URL         string       `json:"Url"`
}

type bingMetadata struct {
	URI  string `json:"Uri"`
	Type string `json:"Type"`
}

const azureURL = "https://api.datamarket.azure.com"

func findBingSearchPath(key string) (string, error) {
	paths := []string{"/Data.ashx/Bing/Search/v1/Web", "/Data.ashx/Bing/SearchWeb/v1/Web"}
	query := "?Query=%27I<3BSW%27"
	for _, path := range paths {
		fullURL := azureURL + path + query
		client := &http.Client{}
		req, err := http.NewRequest("GET", fullURL, nil)
		if err != nil {
			return "", err
		}
		req.SetBasicAuth(key, key)
		resp, err := client.Do(req)
		if err != nil {
			return "", err
		}
		if resp.StatusCode == 200 {
			return path, nil
		}
	}
	return "", errors.New("invalid Bing API key")
}

func bingHTML(domain string) ([]string, error) {
	results := []string{}
	resp, err := http.Get("http://www.bing.com/search?q=site:pastebin.com%20" + domain)
	if err != nil {
		return results, err
	}
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return results, err
	}
	doc.Selection.Find("cite").Each(func(_ int, s *goquery.Selection) {
		results = append(results, s.Text())
	})
	return results, nil
}

func bingAPI(domain, key string) ([]string, error) {
	results := []string{}
	client := &http.Client{}
	path, err := findBingSearchPath(key)
	if err != nil {
		return results, err
	}
	req, err := http.NewRequest("GET", azureURL+path+"?Query=%27site:pastebin.com%20"+domain+"%27&$top=50&Adult=%27off%27&$format=json", nil)
	if err != nil {
		return results, err
	}
	req.SetBasicAuth(key, key)
	resp, err := client.Do(req)
	if err != nil {
		return results, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return results, err
	}
	m := &bingMessage{}
	if err = json.Unmarshal(body, &m); err != nil {
		return results, err
	}
	for _, res := range m.D.Results {
		results = append(results, res.URL)
	}
	return results, nil
}

func main() {

	domain := flag.String("d", "", "domain to search for")
	apiKey := flag.String("k", "", "optional bing api key")
	flag.Parse()

	if *domain == "" {
		log.Fatal("-d required")
	}

	var results []string
	var err error

	if *apiKey != "" {
		results, err = bingAPI(*domain, *apiKey)
		if err != nil {
			log.Fatalf("Error using Bing API. Error %s", err.Error())
		}
	} else {
		results, err = bingHTML(*domain)
		if err != nil {
			log.Fatalf("Error searching Bing. Error %s", err.Error())
		}
	}
	for _, r := range results {
		fmt.Println(r)
	}
}
