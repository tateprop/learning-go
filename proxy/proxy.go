package main

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func parseGetRequest(request *http.Request) *http.Request {
	req := request

	URL := strings.Clone(req.RequestURI)
	req.RequestURI = ""

	u, err := url.Parse(URL)
	if err != nil {
		panic(err)
	}
	req.URL = u

	return req
}
func fowardRequest(request *http.Request) string {
	req := parseGetRequest(request)
	client := &http.Client{}
	response, error := client.Do(req)

	if error != nil {
		fmt.Println(error)
	}

	if response.Header.Get("Content-Encoding") == "gzip" {
		response.Body, error = gzip.NewReader(response.Body)
		if error != nil {
			panic(error)
		}
	}
	responseBody, error := io.ReadAll(response.Body)

	if error != nil {
		fmt.Println(error)
	}
	return string(responseBody)
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	page := fowardRequest(r)
	io.WriteString(w, page)
}

func main() {
	http.HandleFunc("/", getRoot)

	err := http.ListenAndServe(":3333", nil)

	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
