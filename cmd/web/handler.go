package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
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
		log.Print(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < -1 {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Display specific snippet with ID = %d", id)
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)

		// w.WriteHeader(http.StatusMethodNotAllowed)
		// w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
		// or instead
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Create a new snippet..."))
}
