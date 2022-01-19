package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/", HelloServer)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Print(err)
	}
}

const url1 string = "http://httpbin.org/get"
const url2 string = "http://httpbin-svc-external/get"

func HelloServer(w http.ResponseWriter, r *http.Request) {

	if r.Header != nil {
		for k, v := range r.Header {
			fmt.Println(k, v)
		}
	}

	req, err := http.NewRequest("GET", url1, nil)
	if err != nil {
		fmt.Fprintf(w, "Error: %v\n", err)
		return
	}

	addHeader := func(dst *http.Request, src *http.Request, h string) {
		if src.Header.Get(h) != "" {
			dst.Header.Add(h, strings.Join(src.Header.Values(h), ","))
		}
	}

	addHeader(req, r, "x-request-id")
	addHeader(req, r, "x-b3-traceid")
	addHeader(req, r, "x-b3-spanid")
	addHeader(req, r, "x-b3-parentspanid")
	addHeader(req, r, "x-b3-sampled")
	addHeader(req, r, "x-b3-flags")
	addHeader(req, r, "x-ot-span-context")

	http.DefaultClient.Do(req)

	req, err = http.NewRequest("GET", url2, nil)
	if err != nil {
		fmt.Fprintf(w, "Error: %v\n", err)
		return
	}

	addHeader(req, r, "x-request-id")
	addHeader(req, r, "x-b3-traceid")
	addHeader(req, r, "x-b3-spanid")
	addHeader(req, r, "x-b3-parentspanid")
	addHeader(req, r, "x-b3-sampled")
	addHeader(req, r, "x-b3-flags")
	addHeader(req, r, "x-ot-span-context")

	http.DefaultClient.Do(req)

	data, err := ioutil.ReadFile("config/config-data.json")
	if err != nil {
		data = []byte(err.Error())
	}

	fmt.Fprintf(w, "<data from golang-app-data: [%s] >\n", string(data))
}
