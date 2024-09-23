package handlers
import (
	"encoding/json"
	"net/http"
	"strings"
	"html/template"
	"groupie-tracker/internal/models"
	"groupie-tracker/internal/services"
)
func GetAllData(w http.ResponseWriter, r *http.Request) {
	// Create a channel for fetching data
	dataChan := make(chan *models.Datas)
	errChan := make(chan error)
	// Fetch data in a goroutine
	go func() {
		data, err := services.FetchData()
		if err != nil {
			errChan <- err
			return
		}
		dataChan <- data
	}()
	// Wait for data or an error
	select {
	case data := <-dataChan:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(data)
	case <-errChan:
		ErrorHandler(w, r, http.StatusInternalServerError, "Failed to fetch data",)
	}
}
func SearchArtists(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		ErrorHandler(w, r, http.StatusBadRequest, "Search query is required")
		return
	}
	// Create channels for fetching data and filtering artists
	dataChan := make(chan *models.Datas)
	errChan := make(chan error)
	filteredArtistsChan := make(chan []models.Artist)
	// Fetch data in a goroutine
	go func() {
		data, err := services.FetchData()
		if err != nil {
			errChan <- err
			return
		}
		dataChan <- data
	}()
	// Filter artists in another goroutine
	go func() {
		data := <-dataChan
		var filteredArtists []models.Artist
		for _, artist := range data.ArtistsData {
			if strings.Contains(strings.ToLower(artist.Name), strings.ToLower(query)) {
				filteredArtists = append(filteredArtists, artist)
			}
		}
		filteredArtistsChan <- filteredArtists
	}()
	// Wait for filtered artists or an error
	select {
	case filteredArtists := <-filteredArtistsChan:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(filteredArtists)
	case <-errChan:
		ErrorHandler(w, r, http.StatusInternalServerError, "Failed to fetch data")
	}
}
func ErrorHandler(w http.ResponseWriter, r *http.Request, statusCode int, message string) {
	// Parse the HTML template from a file
	tmplPath := "templates/error.html"
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// Define data to pass to the template
	data := struct {
		StatusCode  int
		Message     string
		Description string
	}{
		StatusCode:  statusCode,
		Message:     message,
		Description: http.StatusText(statusCode),
	}
	// Set the status code and render the template
	w.WriteHeader(statusCode)
	tmpl.Execute(w, data)
}
