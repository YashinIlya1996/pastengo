package main

import (
	"errors"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/Yashin1996/pastengo/internal/models"
)

const (
	projFolder    = "/home/metallurg/GolandProjects/pastengo/"
	htmlFolder    = projFolder + "ui/html/"
	pagesFolder   = htmlFolder + "pages/"
	partialFolder = htmlFolder + "partial/"
)

type noDirFS struct {
	fs http.FileSystem
}

func (ndfs noDirFS) Open(name string) (http.File, error) {
	f, err := ndfs.fs.Open(name)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if err != nil {
		closeErr := f.Close()
		if closeErr != nil {
			return nil, closeErr
		}
		return nil, err
	}

	if s.IsDir() {
		index := filepath.Join(name, "index.html")
		if _, err := ndfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}
			return nil, err
		}
	}

	return f, nil
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	files := []string{
		htmlFolder + "base.html",
		pagesFolder + "home.html",
		partialFolder + "nav.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, nil)
		return
	}

	data := &templateDataStruct{
		Snippets: snippets,
	}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, err)
	}

}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < -1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
			return
		}
		app.serverError(w, err)
		return
	}

	files := []string{
		htmlFolder + "base.html",
		partialFolder + "nav.html",
		pagesFolder + "view.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	templateData := &templateDataStruct{
		Snippet: snippet,
	}

	err = ts.ExecuteTemplate(w, "base", templateData)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)

		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Create a new snippet..."))
}
