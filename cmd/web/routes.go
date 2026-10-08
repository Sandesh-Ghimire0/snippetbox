package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", app.home)               // subtree path
	mux.HandleFunc("GET /snippet/{id}", app.showSnippet) // fixed path
	mux.HandleFunc("GET /snippet/create", app.createSnippetForm)
	mux.HandleFunc("POST /snippet/create", app.createSnippet)
	// mux.HandleFunc("/snippet/create/", createSnippet) // subtree path --> if the url contains /snippet/create/** prefix then call createsnippet handler

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
