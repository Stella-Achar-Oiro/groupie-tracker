package handlers

import (
	"fmt"
	"groupie-tracker/api"
	"groupie-tracker/models"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

// ArtistHandler handles requests to display artist information
func ArtistHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the template file
	tmpl, err := template.ParseFiles("templates/artists.html")
	if err != nil {
		log.Fatal(err)
	}

	// Fetch data from APIs
	artists, err := api.GetArtists()
	if err != nil {
		log.Fatal(err)
	}

	relations, err := api.GetRelations()
	if err != nil {
		log.Fatal(err)
	}

	// Get the ID of the Artist from the URL
	artistIDStr := r.URL.Path[9:]
	artistID, err := strconv.Atoi(artistIDStr)
	if err != nil {
		log.Fatal(err)
	}

	// Check if the artistID is valid
	if artistID < 1 || artistID > len(artists) {
		http.Error(w, "Artist not found", http.StatusNotFound)
		return
	}

	// Prepare artist info for the template
	artist := artists[artistID-1]
	artistInfo := map[string]interface{}{
		"Image":        artist.Image,
		"Name":         artist.Name,
		"Members":      artist.Members,
		"CreationDate": artist.CreationDate,
		"FirstAlbum":   artist.FirstAlbum,
		"Locations":    relations.Index[artistID-1].DatesLocations, // Assuming that this is how relations are structured
	}

	// Execute the template with the artistInfo map
	if err := tmpl.Execute(w, artistInfo); err != nil {
		log.Fatal(err)
	}
}

// FilterHandler handles requests to filter artists based on user-provided criteria
func FilterHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the template file
	tmpl, err := template.ParseFiles("templates/filter.html")
	if err != nil {
		log.Fatal(err)
	}

	// Get filter criteria from the request
	filters := api.ExtractFilters(r)

	// Fetch data from APIs
	artists, err := api.GetArtists()
	if err != nil {
		log.Fatal(err)
	}

	locations, err := api.GetLocations()
	if err != nil {
		log.Fatal(err)
	}

	// Map locations.Index to the appropriate structure
	var locationData []models.Location
	for _, loc := range locations.Index {
		locationData = append(locationData, models.Location{
			Index: []struct {
				ID        int      `json:"id"`
				Locations []string `json:"locations"`
				Dates     string   `json:"dates"`
			}{
				{
					ID:        loc.ID,
					Locations: loc.Locations,
					Dates:     loc.Dates,
				},
			},
		})
	}

	// Create a data structure to hold the fetched data
	var data models.Datas
	data.ArtistsData = artists
	data.LocationsData = locationData

	// Check if any filters are applied
	if filters["creationDate"] == "" && filters["firstAlbum"] == "" && filters["members"] == "" && filters["CitySearch"] == "" {
		fmt.Println("No filter applied")
		// Render the template with all artists data
		if err := tmpl.Execute(w, data.ArtistsData); err != nil {
			log.Fatal(err)
		}
	} else {
		// Apply filters to the data
		filteredArtists := api.FilterArtists(filters, data)

		if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
			// If it's an AJAX request, render only the filtered artists
			tmpl, _ := template.New("filtered").Parse(`
					{{range .}}
					<div class="artists">
						<a href='/artists/{{.ID}}'>
							<div class="descrip">
								<h3>{{.Name}}</h3>
							</div>
							<img src="{{.Image}}" alt="{{.Name}}">
						</a>
					</div>
					{{end}}
				`)
			tmpl.Execute(w, filteredArtists)
		} else {
			// For non-AJAX requests, render the full page
			tmpl, _ := template.ParseFiles("templates/filter.html")
			tmpl.Execute(w, filteredArtists)
		}
	}
}

// HomeHandler handles requests to display the home page with artist and location data
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the template file
	tmpl, err := template.ParseFiles("templates/home.html")
	if err != nil {
		log.Fatal(err)
	}

	// Fetch data from APIs
	artists, err := api.GetArtists()
	if err != nil {
		log.Fatal(err)
	}

	locations, err := api.GetLocations()
	if err != nil {
		log.Fatal(err)
	}

	// Map locations.Index to the appropriate structure
	var locationData []models.Location
	for _, loc := range locations.Index {
		locationData = append(locationData, models.Location{
			Index: []struct {
				ID        int      `json:"id"`
				Locations []string `json:"locations"`
				Dates     string   `json:"dates"`
			}{
				{
					ID:        loc.ID,
					Locations: loc.Locations,
					Dates:     loc.Dates,
				},
			},
		})
	}

	// Create a data structure to hold the fetched data
	data := models.Datas{
		ArtistsData:   artists,
		LocationsData: locationData, // Now contains the correctly mapped data
	}

	// Render the template with the data
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// IndexHandler serves the index page template
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the index.html template
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		log.Fatal("Error parsing template:", err)
	}

	// Render the template
	if err := tmpl.Execute(w, nil); err != nil {
		log.Fatal("Error executing template:", err)
	}
}
