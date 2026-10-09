package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// session is only required in dynamic routes, /static/ does not need session
	dynamicMiddleWare := alice.New(app.session.LoadAndSave)
	mux := http.NewServeMux()

	mux.Handle("GET /", dynamicMiddleWare.ThenFunc(app.home))                    // subtree path
	mux.Handle("GET /snippet/{id}", dynamicMiddleWare.ThenFunc(app.showSnippet)) // fixed path
	mux.Handle("POST /snippet/create", dynamicMiddleWare.ThenFunc(app.createSnippet))
	mux.Handle("GET /snippet/create", dynamicMiddleWare.ThenFunc(app.createSnippetForm))
	// mux.HandleFunc("/snippet/create/", createSnippet) // subtree path --> if the url contains /snippet/create/** prefix then call createsnippet handler

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
