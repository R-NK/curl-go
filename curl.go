package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/jessevdk/go-flags"
)

func isValidURL(u string) bool {
	_, err := url.ParseRequestURI(u)
	if err != nil {
		return false
	}
	return true
}

func isValidMethod(method string) bool {
	switch method {
	case "GET":
		return true
	case "POST":
		return true
	default:
		return false
	}
}

func parseHeader(rawHeader string) (key string, value string) {
	headerPair := strings.Split(rawHeader, ":")
	key = headerPair[0]
	value = strings.TrimLeft(headerPair[1], " ")
	return key, value
}

// Get method
func Get(u string, headers []string) (*http.Response, error) {
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	if len(headers) > 0 {
		for _, header := range headers {
			headerKey, headerValue := parseHeader(header)
			req.Header.Set(headerKey, headerValue)
		}
	}
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	return resp, nil
}

// Post method
func Post(u string, headers []string, body string) (*http.Response, error) {
	req, err := http.NewRequest("POST", u, strings.NewReader(body))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	if len(headers) > 0 {
		for _, header := range headers {
			headerKey, headerValue := parseHeader(header)
			req.Header.Set(headerKey, headerValue)
		}
	}
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	return resp, nil
}

func main() {
	var opts struct {
		Method string   `short:"X" description:"HTTP Method" default:"GET"`
		Header []string `short:"H" long:"header" description:"Change HTTP Header" default:"Content-Type:application/x-www-form-urlencoded"`
		Data   string   `short:"d" long:"data" description:"Sends the specified data in a POST request to the HTTP server"`
	}
	args, err := flags.Parse(&opts)
	if err != nil {
		panic(err)
	}
	url := args[0]
	if !isValidURL(url) {
		log.Fatal("invalid URL!")
		os.Exit(1)
	}

	if !isValidMethod(opts.Method) {
		log.Fatal("invalid method!")
		os.Exit(1)
	}

	var resp *http.Response
	switch opts.Method {
	case "GET":
		resp, err = Get(url, opts.Header)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
	case "POST":
		resp, err = Post(url, opts.Header, opts.Data)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
	}

	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("%s", string(b))
}
