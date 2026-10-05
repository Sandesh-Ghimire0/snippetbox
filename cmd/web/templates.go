package main

import "github.com/Sandesh-Ghimire0/snippetbox/pkg/models"

/*
Go’s html/template package allows you to pass in one — and only one — item of dynamic data when rendering a template.
But in a real-world application there are often multiple pieces of dynamic data that you want to display in the same page.

A lightweight and type-safe way to achieve this is to wrap your dynamic data in a struct which acts like a
single ‘holding structure’ for your data.
*/


// holing structure for any kind of dynamic data that we want to pass in our HTML templates.
type templateData struct {
	Snippet *models.Snippet
}
