package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

var portnumber = 8087

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/hello" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Only GET method is supported.", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "Hello Robert!")
}

func startForm() {

	fileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileServer)
	http.HandleFunc("/form", formHandler)
	http.HandleFunc("/hello", helloHandler)

	fmt.Printf("Starting server at port %s\n", port)
	if err := http.ListenAndServe(string(":"+strconv.Itoa(portnumber)), nil); err != nil {
		log.Fatal(err)
	}
}

func formHandler(w http.ResponseWriter, r *http.Request) {

	createPost(w, r)
	//	router.HandleFunc("/posts", createPost).Methods("POST")

	/*
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		fmt.Fprintf(w, "POST request successful\n")
		name := r.FormValue("name")
		address := r.FormValue("address")

		fmt.Fprintf(w, "Name = %s\n", name)
		fmt.Fprintf(w, "Address = %s\n", address)
		//*/
}
