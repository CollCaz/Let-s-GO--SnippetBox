package main

import (
	"path/filepath"
	"text/template"

	"github.com/CollCaz/Lets-GO--SnippetBox/internal/models"
)

type templateData struct {
	Snippet  *models.Snippet   // For displaying a single Snippet
	Snippets []*models.Snippet // For displaying a bunch of Snippets
}

// Creates in memory cache for all templates
func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// all template pages
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		// first, parse the base template into a template set
		ts, err := template.ParseFiles("./ui/html/base.tmpl.html")
		if err != nil {
			return nil, err
		}

		// Then parse all the partials on the first tempalte set to add any partials
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl.html")
		if err != nil {
			return nil, err
		}

		// Finally parse the page into this template set
		ts, err = ts.ParseGlob(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts

	}

	return cache, nil
}
