package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"path/filepath"
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
	const (
		projFolder = "/home/metallurg/GolandProjects/pastengo/"
		htmlFolder  = projFolder + "ui/html/"
		pagesFolder = htmlFolder + "pages/"
		partialFolder = htmlFolder + "partial/"
	)
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

	err = ts.ExecuteTemplate(w, "base", nil)
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
	fmt.Fprintf(w, "Display specific snippet with ID = %d", id)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)

		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Create a new snippet..."))
}
