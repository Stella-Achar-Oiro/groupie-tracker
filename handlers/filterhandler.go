package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	files "groupie-tracker/file"
	"groupie-tracker/models"
)

func FilterHandler(w http.ResponseWriter, r *http.Request) {
	file, error := template.ParseFiles("./templates/filter.html")

	if error != nil {
		log.Fatal(error)
	}

	// Get filters
	filters := files.GetFilters(r)

	var data models.Datas

	// Open JSON files
	data = files.OpenJSON("artists.json", data)
	data = files.OpenJSON("locations.json", data)

	// check if there are datas into filters on the HTMLpage
	if filters["creationDate"] == "" && filters["firstAlbum"] == "" && filters["members"] == "" && filters["CitySearch"] == "" {
		fmt.Println("No filter")
		// if there isn't filters, we execute the templates with the data from Json files (render like mainPage)
		file.Execute(w, data.ArtistsData)
	} else {
		// if there are filters, we call the function dataToPush and execute the template with filter
		toPush := files.DataToPush(filters, data)

		// Render template with artist matching filters
		file.Execute(w, toPush)
	}
}
