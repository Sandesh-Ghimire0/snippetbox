package main

import (
	"html/template"
	"net/url"
	"path"
	"path/filepath"
	"time"

	"github.com/Sandesh-Ghimire0/snippetbox/pkg/models"
)

/*
Go’s html/template package allows you to pass in one — and only one — item of dynamic data when rendering a template.
But in a real-world application there are often multiple pieces of dynamic data that you want to display in the same page.

A lightweight and type-safe way to achieve this is to wrap your dynamic data in a struct which acts like a
single ‘holding structure’ for your data.
*/

// holing structure for any kind of dynamic data that we want to pass in our HTML templates.
type templateData struct {
	CurrentYear int
	FormData    url.Values
	FormErrors  map[string]string
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
}

// custom template function
// converts the time into human date format
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// object containing custom template function
var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {
	/*
		1. Create an empty cache map
		2. Get all the pages like (show.page.tmpl)
		3. Iterate through the page and do following:
			- parse the page and get the template set
			- Get the layout and partial templates
			- add parse layout and partial to the ts
			- add page in the cache map
	*/

	cache := map[string]*template.Template{}

	// slice of relative paths to current working directory
	pages, err := filepath.Glob(path.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		// The template.FuncMap must be registered with the template set before you
		// call the ParseFiles() method. This means we have to use template.New() to
		// create an empty template set, use the Funcs() method to register the
		// template.FuncMap, and then parse the file as normal.

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// get layout tmpl and add to ts
		ts, err = ts.ParseGlob(path.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		// get partial tmpl and add to ts
		ts, err = ts.ParseGlob(path.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil

}
