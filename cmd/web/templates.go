package main

import (
	"net/http"
	"path/filepath"
	"text/template"
	"time"

	"github.com/CollCaz/Lets-GO--SnippetBox/internal/models"
)

type templateData struct {
	CurrentYear int
	Snippet     *models.Snippet   // For displaying a single Snippet
	Snippets    []*models.Snippet // For displaying a bunch of Snippets
}

func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
	}
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
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
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/pages/base.tmpl.html")
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
