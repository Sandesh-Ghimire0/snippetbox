package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

func (app *application) ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

	// return internal server error as response
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// The clientError helper sends a specific status code and corresponding description
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) NotFound(w http.ResponseWriter) {
	// not found is also a client error
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {

	ts, ok := app.templateCache[name]
	if !ok {
		app.ServerError(w, fmt.Errorf("The template %s does not exist", name))
		return
	}

	// catching the runtime template error (runtime error might return the half rendered html with error)
	// at first render the template in the buffer
	// if template executes successfully then write the http response
	// else return the error

	buf := new(bytes.Buffer)

	err := ts.Execute(buf, app.addDefaultData(td, r))
	if err != nil {
		app.ServerError(w, err)
		return
	}

	buf.WriteTo(w)
}

// this will add the default data to the templates
func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}
	td.CurrentYear = time.Now().Year()
	return td
}
