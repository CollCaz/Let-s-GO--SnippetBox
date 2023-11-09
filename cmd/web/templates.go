package main

import (
	"io/fs"
	"net/http"
	"path/filepath"
	"text/template"
	"time"

	"github.com/CollCaz/Lets-GO--SnippetBox/internal/models"
	"github.com/CollCaz/Lets-GO--SnippetBox/ui"
	"github.com/justinas/nosurf"
)

type templateData struct {
	CurrentYear     int
	Snippet         *models.Snippet   // For displaying a single Snippet
	Snippets        []*models.Snippet // For displaying a bunch of Snippets
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear:     time.Now().Year(),
		Flash:           app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: app.isAuthenticated(r),
		CSRFToken:       nosurf.Token(r),
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
	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		patterns := []string{
			"html/pages/base.tmpl.html",
			"html/partials/*.tmpl.html",
			page,
		}

		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}
