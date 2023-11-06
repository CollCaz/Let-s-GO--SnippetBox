package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"

	"github.com/CollCaz/Lets-GO--SnippetBox/internal/models"
)

// The handler func for /
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/pages/base.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
	}
	template_set, err := template.ParseFiles(files...)
	if err != nil {
		log.Print(err.Error())
		app.serverError(w, err)
	}
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	err = template_set.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Print(err.Error())
		app.serverError(w, err)
	}
}

// Handler func for Viewing snippets with specific ID
// snippetbox.com/snippet/view?id=123
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrorNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	fmt.Fprintf(w, "%+v", snippet)
}

// Handler func for creating snippets, only accepts POST requests
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	title := "Emilia<3"
	content := "Emilia is best girl"
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}
