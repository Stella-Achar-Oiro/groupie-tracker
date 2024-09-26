package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	ArtistsAPI   = "https://groupietrackers.herokuapp.com/api/artists"
	LocationsAPI = "https://groupietrackers.herokuapp.com/api/locations"
	DatesAPI     = "https://groupietrackers.herokuapp.com/api/dates"
	RelationsAPI = "https://groupietrackers.herokuapp.com/api/relation"
)

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
		ID             int                 `json:"id"`
		DatesLocations map[string][]string `json:"datesLocations"`
	} `json:"index"`
}

type Datas struct {
	ArtistsData   []Artist `json:"artists"`
	LocationsData Location `json:"locations"`
	DatesData     Date     `json:"dates"`
	RelationsData Relation `json:"relations"`
}

type SearchResult struct {
	Artists []Artist `json:"artists"`
}

type Suggestion struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

type GeoLocation struct {
	Address string  `json:"address"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
}

type ArtistDetail struct {
	Artist    Artist              `json:"artist"`
	Locations []GeoLocation       `json:"locations"`
	Dates     []string            `json:"dates"`
	Relations map[string][]string `json:"relations"`
}

type Cache struct {
	data      Datas
	expiresAt time.Time
	mutex     sync.RWMutex
}

const (
	MapboxAccessToken  = "pk.eyJ1Ijoic3RlbGxhYWNoYXJvaXJvIiwiYSI6ImNtMWhmZHNlODBlc3cybHF5OWh1MDI2dzMifQ.wk3v-v7IuiSiPwyq13qdHw"
	MapboxGeocodingAPI = "https://api.mapbox.com/geocoding/v5/mapbox.places"
)

var (
	cache    Cache
	indexTpl *template.Template
)

const cacheDuration = 1 * time.Hour

func geocode(address string) (GeoLocation, error) {
	url := fmt.Sprintf("%s/%s.json?access_token=%s", MapboxGeocodingAPI, url.QueryEscape(address), MapboxAccessToken)

	resp, err := http.Get(url)
	if err != nil {
		return GeoLocation{}, err
	}
	defer resp.Body.Close()

	var result struct {
		Features []struct {
			Center [2]float64 `json:"center"`
		} `json:"features"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return GeoLocation{}, err
	}

	if len(result.Features) == 0 {
		return GeoLocation{}, fmt.Errorf("no results found for address: %s", address)
	}

	return GeoLocation{
		Address: address,
		Lon:     result.Features[0].Center[0],
		Lat:     result.Features[0].Center[1],
	}, nil
}

func main() {
	// Initialize cache
	cache = Cache{}

	// Initial data fetch
	if err := refreshCache(); err != nil {
		log.Fatal("Failed to fetch initial data:", err)
	}

	// Parse HTML template
	var err error
	indexTpl, err = template.ParseFiles("index.html")
	if err != nil {
		log.Fatal("Failed to parse template:", err)
	}

	// Set up routes
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/api/search", handleSearch)
	http.HandleFunc("/api/artist/", handleArtist)
	http.HandleFunc("/api/suggestions", handleSuggestions)

	// Start server
	fmt.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func refreshCache() error {
	var newData Datas
	err := fetchAllData(&newData)
	if err != nil {
		return err
	}

	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	cache.data = newData
	cache.expiresAt = time.Now().Add(cacheDuration)

	return nil
}

func getCachedData() (Datas, error) {
	cache.mutex.RLock()
	if time.Now().Before(cache.expiresAt) {
		defer cache.mutex.RUnlock()
		return cache.data, nil
	}
	cache.mutex.RUnlock()

	if err := refreshCache(); err != nil {
		return Datas{}, err
	}

	return cache.data, nil
}

func fetchAllData(data *Datas) error {
	var wg sync.WaitGroup
	errChan := make(chan error, 4)

	wg.Add(4)
	go fetchData(ArtistsAPI, &data.ArtistsData, &wg, errChan)
	go fetchData(LocationsAPI, &data.LocationsData, &wg, errChan)
	go fetchData(DatesAPI, &data.DatesData, &wg, errChan)
	go fetchData(RelationsAPI, &data.RelationsData, &wg, errChan)

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func fetchData(url string, target interface{}, wg *sync.WaitGroup, errChan chan<- error) {
	defer wg.Done()
	resp, err := http.Get(url)
	if err != nil {
		errChan <- fmt.Errorf("failed to fetch data from %s: %v", url, err)
		return
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		errChan <- fmt.Errorf("failed to decode data from %s: %v", url, err)
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	err := indexTpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
		return
	}
}

func handleSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Missing search query", http.StatusBadRequest)
		return
	}

	cachedData, err := getCachedData()
	if err != nil {
		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
		return
	}

	results := searchArtists(query, cachedData.ArtistsData)
	json.NewEncoder(w).Encode(results)
}

func searchArtists(query string, artists []Artist) SearchResult {
	var results SearchResult
	lowercaseQuery := strings.ToLower(query)

	for _, artist := range artists {
		if strings.Contains(strings.ToLower(artist.Name), lowercaseQuery) ||
			containsAny(artist.Members, lowercaseQuery) ||
			strings.Contains(strings.ToLower(artist.FirstAlbum), lowercaseQuery) ||
			strconv.Itoa(artist.CreationDate) == query {
			results.Artists = append(results.Artists, artist)
		}
	}

	return results
}

func containsAny(slice []string, substr string) bool {
	for _, s := range slice {
		if strings.Contains(strings.ToLower(s), substr) {
			return true
		}
	}
	return false
}

func handleArtist(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/artist/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid artist ID", http.StatusBadRequest)
		return
	}

	cachedData, err := getCachedData()
	if err != nil {
		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
		return
	}

	var artist Artist
	for _, a := range cachedData.ArtistsData {
		if a.ID == id {
			artist = a
			break
		}
	}

	if artist.ID == 0 {
		http.Error(w, "Artist not found", http.StatusNotFound)
		return
	}

	details := ArtistDetail{
		Artist:    artist,
		Locations: getLocations(id, cachedData.LocationsData),
		Dates:     getDates(id, cachedData.DatesData),
		Relations: getRelations(id, cachedData.RelationsData),
	}

	json.NewEncoder(w).Encode(details)
}

func getLocations(id int, locationsData Location) []GeoLocation {
	for _, loc := range locationsData.Index {
		if loc.ID == id {
			var geoLocations []GeoLocation
			for _, location := range loc.Locations {
				geoLoc, err := geocode(location)
				if err != nil {
					log.Printf("Failed to geocode location: %v", err)
					continue
				}
				geoLocations = append(geoLocations, geoLoc)
			}
			return geoLocations
		}
	}
	return nil
}


func getDates(id int, datesData Date) []string {
	for _, date := range datesData.Index {
		if date.ID == id {
			return date.Dates
		}
	}
	return nil
}

func getRelations(id int, relationsData Relation) map[string][]string {
	for _, rel := range relationsData.Index {
		if rel.ID == id {
			return rel.DatesLocations
		}
	}
	return nil
}

func handleSuggestions(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Missing search query", http.StatusBadRequest)
		return
	}

	cachedData, err := getCachedData()
	if err != nil {
		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
		return
	}

	suggestions := getSuggestions(query, cachedData.ArtistsData)
	json.NewEncoder(w).Encode(suggestions)
}

func getSuggestions(query string, artists []Artist) []Suggestion {
	var suggestions []Suggestion
	lowercaseQuery := strings.ToLower(query)

	for _, artist := range artists {
		if strings.Contains(strings.ToLower(artist.Name), lowercaseQuery) {
			suggestions = append(suggestions, Suggestion{Text: artist.Name, Type: "artist"})
		}
		for _, member := range artist.Members {
			if strings.Contains(strings.ToLower(member), lowercaseQuery) {
				suggestions = append(suggestions, Suggestion{Text: member, Type: "member"})
			}
		}
		if strings.Contains(strings.ToLower(artist.FirstAlbum), lowercaseQuery) {
			suggestions = append(suggestions, Suggestion{Text: artist.FirstAlbum, Type: "album"})
		}
		if strings.Contains(strconv.Itoa(artist.CreationDate), query) {
			suggestions = append(suggestions, Suggestion{Text: strconv.Itoa(artist.CreationDate), Type: "year"})
		}
	}

	return suggestions
}