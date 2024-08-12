package handlers

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	files "groupie-tracker/file"
	"groupie-tracker/models"
)

func ArtistHandler(w http.ResponseWriter, r *http.Request) {
	tmp, error := template.ParseFiles("./templates/artists.html")

	// ART var of type DATAS struct
	var ART models.Datas

	// open JSON files
	ART = files.OpenJSON("artists.json", ART)
	ART = files.OpenJSON("relation.json", ART)

	// Get the ID of the Artist
	id := r.URL.Path[9:]
	p, _ := strconv.Atoi(id)

	if error != nil {
		log.Fatal(error)
	}

	infoArt := make(map[string]interface{})

	// Store usefull information in infoArt map
	infoArt["IMAGE"] = ART.ArtistsData[p-1].IMAGE
	infoArt["NAME"] = ART.ArtistsData[p-1].NAME
	infoArt["MEMBERS"] = ART.ArtistsData[p-1].MEMBERS
	infoArt["CREA_DATE"] = ART.ArtistsData[p-1].CREA_DATE
	infoArt["FIRST_ALBUM"] = ART.ArtistsData[p-1].FIRST_ALBUM
	infoArt["DATESLOCAT"] = ART.RelationData[p-1].DATESLOCAT

	if error := tmp.Execute(w, infoArt); error != nil {
		log.Fatal(error)
	}
}
