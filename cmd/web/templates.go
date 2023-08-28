package main

import "github.com/Yashin1996/pastengo/internal/models"

type templateDataStruct struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}
