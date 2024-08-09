package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

// Structs for each API endpoint
type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Locations    string   `json:"locations"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

type Location struct {
	Index []struct {
		ID        int      `json:"id"`
		Locations []string `json:"locations"`
		Dates     string   `json:"dates"`
	} `json:"index"`
}

type Date struct {
	Index []struct {
		ID    int      `json:"id"`
		Dates []string `json:"dates"`
	} `json:"index"`
}

type Relation struct {
	Index []struct {
		ID             int                    `json:"id"`
		DatesLocations map[string]interface{} `json:"datesLocations"`
	} `json:"index"`
}

// Templates
var (
	tmplHome      *template.Template
	tmplArtists   *template.Template
	tmplLocations *template.Template
	tmplDates     *template.Template
	tmplRelation  *template.Template
	once          sync.Once
)

func main() {
	// Initialize templates
	initTemplates()

	mux := http.NewServeMux()
	mux.HandleFunc("/", handleHome)
	mux.HandleFunc("/artists", handleArtists)
	mux.HandleFunc("/locations", handleLocations)
	mux.HandleFunc("/dates", handleDates)
	mux.HandleFunc("/relation", handleRelation)
	http.HandleFunc("/search", handleSearch)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", mux)
}

// Initialize templates only once
func initTemplates() {
	once.Do(func() {
		tmplHome = template.Must(template.ParseFiles("templates/home.html"))
		tmplArtists = template.Must(template.ParseFiles("templates/artists.html"))
		tmplLocations = template.Must(template.ParseFiles("templates/locations.html"))
		tmplDates = template.Must(template.ParseFiles("templates/dates.html"))
		tmplRelation = template.Must(template.ParseFiles("templates/relation.html"))
	})
}

// Handler functions
func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	renderTemplate(w, tmplHome, nil)
}

// func handleArtists(w http.ResponseWriter, r *http.Request) {
// 	data, err := fetchData[[]Artist]("https://groupietrackers.herokuapp.com/api/artists")
// 	if err != nil {
// 		handleError(w, err)
// 		return
// 	}

// 	renderTemplate(w, tmplArtists, data)
// }

const itemsPerPage = 10

func handleArtists(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	artists, err := fetchData[[]Artist]("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		handleError(w, err)
		return
	}

	startIndex := (page - 1) * itemsPerPage
	endIndex := startIndex + itemsPerPage
	if endIndex > len(artists) {
		endIndex = len(artists)
	}

	data := struct {
		Title       string
		Artists     []Artist
		CurrentPage int
		TotalPages  int
		HasPrevPage bool
		HasNextPage bool
	}{
		Title:       "Artists",
		Artists:     artists[startIndex:endIndex],
		CurrentPage: page,
		TotalPages:  (len(artists) + itemsPerPage - 1) / itemsPerPage,
		HasPrevPage: page > 1,
		HasNextPage: endIndex < len(artists),
	}
	renderTemplate(w, tmplArtists, data)
}

func handleLocations(w http.ResponseWriter, r *http.Request) {
	data, err := fetchData[Location]("https://groupietrackers.herokuapp.com/api/locations")
	if err != nil {
		handleError(w, err)
		return
	}

	renderTemplate(w, tmplLocations, data.Index)
}

func handleDates(w http.ResponseWriter, r *http.Request) {
	data, err := fetchData[Date]("https://groupietrackers.herokuapp.com/api/dates")
	if err != nil {
		handleError(w, err)
		return
	}

	renderTemplate(w, tmplDates, data.Index)
}

func handleRelation(w http.ResponseWriter, r *http.Request) {
	data, err := fetchData[Relation]("https://groupietrackers.herokuapp.com/api/relation")
	if err != nil {
		handleError(w, err)
		return
	}

	renderTemplate(w, tmplRelation, data.Index)
}

func handleSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	artists, err := fetchData[[]Artist]("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		handleError(w, err)
		return
	}

	var results []Artist
	for _, artist := range artists {
		if strings.Contains(strings.ToLower(artist.Name), strings.ToLower(query)) {
			results = append(results, artist)
		}
	}

	tmpl, err := template.ParseFiles("templates/artists.html")
	if err != nil {
		handleError(w, err)
		return
	}

	err = tmpl.Execute(w, results)
	if err != nil {
		handleError(w, err)
	}
}

// Helper functions
func fetchData[T any](endpoint string) (T, error) {
	var data T
	resp, err := http.Get(endpoint)
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return data, err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}

func renderTemplate(w http.ResponseWriter, tmpl *template.Template, data interface{}) {
	err := tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
