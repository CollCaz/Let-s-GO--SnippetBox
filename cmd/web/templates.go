package main

import "github.com/CollCaz/Lets-GO--SnippetBox/internal/models"

type templateData struct {
	Snippet  *models.Snippet   // For displaying a single Snippet
	Snippets []*models.Snippet // For displaying a bunch of Snippets
}
