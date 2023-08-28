package main

import (
	"html/template"
	"path/filepath"

	"github.com/Yashin1996/pastengo/internal/models"
)

type templateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(pagesFolder + "*html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		ts, err := template.ParseFiles(htmlFolder + "base.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(partialFolder + "*.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[filepath.Base(page)] = ts
	}

	return cache, nil
}
