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

const url string = "http://springboot-app-svc:8080"

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

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Fprintf(w, "Error: %v\n", err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(w, "Upstream error: %d\n", resp.StatusCode)
		return
	}
	if resp.Body == nil {
		fmt.Fprintf(w, "Upstream error: No content\n")
		return
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(w, "Upstream error: Read body failed\n")
		return
	}

	fmt.Fprintf(w, "Gateway get from upstream: [\n\n %s \n]\n", string(b))
}
