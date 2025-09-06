package main

import "github.com/gitavk/snippetbox/internal/models"

// Define a templateData type to act as the holding structure for
type templateData struct {
	Snippet  models.Snippet
	Snippets []models.Snippet
}
