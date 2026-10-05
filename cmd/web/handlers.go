package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/Sandesh-Ghimire0/snippetbox/pkg/models"
)

// this is the method against the application struct
func (app *application) home(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		app.NotFound(w)
		return // most return otherwise below code will execute
	}

	s, err := app.snippets.Latest()

	if err != nil {
		app.ServerError(w, err)
		return
	}

	data := &templateData{Snippets: s}
	files := []string{
		"./ui/html/home.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	// Use the ExecuteTemplate() method to write the content of the "base"
	// template as the response body.
	err = ts.Execute(w, data)
	if err != nil {
		app.ServerError(w, err)
	}
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.NotFound(w)
		return
	}

	s, err := app.snippets.Get(id)

	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.NotFound(w)
		} else {
			app.ServerError(w, err)
		}

		return

	}

	data := &templateData{Snippet: s}
	files := []string{
		"./ui/html/show.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	err = ts.Execute(w, data)
	if err != nil {
		app.ServerError(w, err)
		return
	}

}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		// most set before writeHeader() or write() otherwise it will have no effect
		w.Header().Set("Allow", http.MethodPost)
		// http.Error(w, "Method Not Allowed", 405)
		app.clientError(w, http.StatusMethodNotAllowed)
		// w.WriteHeader(405)
		// w.Write([]byte("Method Not Allowed"))
		return
	}

	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly,slowly!\n\n- Kobayashi Issa"
	expires := "7"

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	// redirect user to the relavant page
	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)

}
