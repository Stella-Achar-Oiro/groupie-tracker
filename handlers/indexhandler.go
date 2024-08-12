package handlers

import (
	"html/template"
	"log"
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	// Set index.html as template
	file, error := template.ParseFiles("./templates/index.html")

	if error != nil {
		log.Fatal(error)
	}

	// Render template
	if error := file.Execute(w, file); error != nil {
		log.Fatal(error)
	}
}
