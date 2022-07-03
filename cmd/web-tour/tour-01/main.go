package main

import (
	"fmt"
	"log"
	"net/http"
)

func init() { log.SetFlags(log.LstdFlags) }

func index(w http.ResponseWriter, req *http.Request) {
	log.Printf("URL.Path = %q", req.URL.Path)
	fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
}

func echo(w http.ResponseWriter, req *http.Request) {
	log.Printf("Headers: %v", req.Header)
	for k, v := range req.Header {
		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
	}
}

// func router() {
// 	http.HandleFunc("/", index)
// 	http.HandleFunc("/echo", echo)
// }

func run(port int) {
	http.HandleFunc("/", index)
	http.HandleFunc("/echo", echo)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func main() {
	run(9001)
}
