package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

// dependency injection
// defining application struct to hold application wide dependencies for the web application
type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP Network Address")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	mux.HandleFunc("/", app.home)               // subtree path
	mux.HandleFunc("/snippet", app.showSnippet) // fixed path
	mux.HandleFunc("/snippet/create", app.createSnippet)
	// mux.HandleFunc("/snippet/create/", createSnippet) // subtree path --> if the url contains /snippet/create/** prefix then call createsnippet handler

	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// By default go http server logs the error to the standard logger if we want to implement our errorLog to log error
	// we need to create a http.Server struct and add the required fields

	server := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}

	infoLog.Printf("starting the server on %s\n", *addr)
	err := server.ListenAndServe()

	errorLog.Fatal(err)

}
