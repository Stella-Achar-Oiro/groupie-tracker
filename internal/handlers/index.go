package handlers

import (
	"html/template"
	"net/http"
)

// HandleIndex returns a http.HandlerFunc that renders the index template
func HandleIndex(tpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := tpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
			return
		}
	}
}