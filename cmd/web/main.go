package main

import (
	"log"
	"net/http"
)


func main() {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.HandleFunc("/", home)               // subtree path
	mux.HandleFunc("/snippet", showSnippet) // fixed path
	mux.HandleFunc("/snippet/create", createSnippet)
	// mux.HandleFunc("/snippet/create/", createSnippet) // subtree path --> if the url contains /snippet/create/** prefix then call createsnippet handler


	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	log.Println("starting the server on :4000")
	err := http.ListenAndServe(":4000", mux)

	log.Fatal(err)

}
