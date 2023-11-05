package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// The HTTP Address
	addr := flag.String("addr", ":4000", "HTTP Network Address")
	flag.Parse()

	mux := http.NewServeMux()

	fileserver := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("/static/", http.StripPrefix("/static", fileserver))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	fmt.Printf("http://www.localhost%s\n", *addr)
	log.Printf("Starting server on %s\n", *addr)
	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}
