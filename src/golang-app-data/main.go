package main

import (
	"fmt"
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

const url string = "http://httpbin.org/get"

func HelloServer(w http.ResponseWriter, r *http.Request) {

	if r.Header != nil {
		for k, v := range r.Header {
			fmt.Println(k, v)
		}
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Fprintf(w, "Error: %v\n", err)
		return
	}

	addHeader := func(dst *http.Request, src *http.Request, h string) {
		if src.Header.Get(h) != "" {
			dst.Header.Add(h, strings.Join(src.Header.Values(h), ","))
		}
	}

	{
		addHeader(req, r, "x-request-id")
		addHeader(req, r, "x-b3-traceid")
		addHeader(req, r, "x-b3-spanid")
		addHeader(req, r, "x-b3-parentspanid")
		addHeader(req, r, "x-b3-sampled")
		addHeader(req, r, "x-b3-flags")
		addHeader(req, r, "x-ot-span-context")

		addHeader(req, r, "x-cloud-trace-context")
		addHeader(req, r, "traceparent")
		addHeader(req, r, "grpc-trace-bin")
	}

	http.DefaultClient.Do(req)
	fmt.Fprintf(w, "<data from golang-app-data>\n")
}
