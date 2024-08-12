package handlers

import (
	"html/template"
	"log"
	"net/http"

	files "groupie-tracker/file"
	"groupie-tracker/models"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	file, error := template.ParseFiles("./templates/home.html")

	if error != nil {
		log.Fatal(error)
	}

	// data var of type DATAS struct
	var data models.Datas

	// Open JSON files
	data = files.OpenJSON("artists.json", data)

	data = files.OpenJSON("locations.json", data)

	// Render template with value stored on data variable
	file.Execute(w, data)
}
