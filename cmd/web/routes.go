package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	// session is only required in dynamic routes, /static/ does not need session
	dynamicMiddleWare := alice.New(app.session.LoadAndSave, noSurf, app.authenticate)
	mux := http.NewServeMux()

	mux.Handle("GET /", dynamicMiddleWare.ThenFunc(app.home))                    // subtree path
	mux.Handle("GET /snippet/{id}", dynamicMiddleWare.ThenFunc(app.showSnippet)) // fixed path
	mux.Handle("POST /snippet/create", dynamicMiddleWare.Append(app.requiresAuthentication).ThenFunc(app.createSnippet))
	mux.Handle("GET /snippet/create", dynamicMiddleWare.Append(app.requiresAuthentication).ThenFunc(app.createSnippetForm))
	// mux.HandleFunc("/snippet/create/", createSnippet) // subtree path --> if the url contains /snippet/create/** prefix then call createsnippet handler

	mux.Handle("GET /users/signup", dynamicMiddleWare.ThenFunc(app.signupUserForm))
	mux.Handle("POST /users/signup", dynamicMiddleWare.ThenFunc(app.sigupUser))
	mux.Handle("GET /users/login", dynamicMiddleWare.ThenFunc(app.loginUserForm))
	mux.Handle("POST /users/login", dynamicMiddleWare.ThenFunc(app.loginUser))
	mux.Handle("POST /users/logout", dynamicMiddleWare.Append(app.requiresAuthentication).ThenFunc(app.logoutUser))

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
