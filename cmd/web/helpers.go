package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/go-playground/form"
)

// the serverError helper writes an error message and stack trace to the errorLog.
// then sends a generic 500 internal server error to the user
func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())

	// We will ignore this error becasue we should decided
	// what to do with it on a case by case basis
	// I don't like this but thats how the book does it
	// Should revisit this later when im better at Go
	app.errorLog.Output(3, trace) //nolint:all FIXME

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the tempalte %s does not exist", page)
		app.serverError(w, err)
		return
	}

	// Initilize a new buffer to render the page two
	buf := new(bytes.Buffer)

	// Write the page to the buffer and check for errors
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(status)
	_, err = buf.WriteTo(w)
	if err != nil {
		app.serverError(w, err)
		return
	}
}

func (app *application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		var invalidDecodeError *form.InvalidDecoderError

		if errors.As(err, &invalidDecodeError) {
			panic(err)
		}
		return err
	}
	return nil
}

func (app *application) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}

	return isAuthenticated
}
