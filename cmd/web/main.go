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

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	// By default go http server logs the error to the standard logger if we want to implement our errorLog to log error
	// we need to create a http.Server struct and add the required fields

	server := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("starting the server on %s\n", *addr)
	err := server.ListenAndServe()

	errorLog.Fatal(err)

}
