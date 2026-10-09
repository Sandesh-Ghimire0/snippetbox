package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Sandesh-Ghimire0/snippetbox/pkg/models/mysql"
	"github.com/alexedwards/scs/v2"
	_ "github.com/go-sql-driver/mysql"
)

// dependency injection
// defining application struct to hold application wide dependencies for the web application
type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *mysql.SnippetModel
	session       *scs.SessionManager
	templateCache map[string]*template.Template
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP Network Address")
	dsn := flag.String("dsn", "web:sanghi@/snippetbox?parseTime=true", "MySQL Data Source Name")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// parses all the pages during start and stores in cache
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	sessionManager := scs.New()
	sessionManager.Lifetime = 12 * time.Hour

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &mysql.SnippetModel{DB: db},
		session:       sessionManager,
		templateCache: templateCache,
	}

	// By default go http server logs the error to the standard logger if we want to implement our errorLog to log error
	// we need to create a http.Server struct and add the required fields

	server := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("starting the server on %s\n", *addr)
	err = server.ListenAndServe()

	errorLog.Fatal(err)

}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
