package api

import (
	"encoding/json"
	"io"
	"net/http"

	"groupie-tracker/models"
)

// FetchData fetches data from a given URL and unmarshals it into the provided data structure
func FetchData(url string, target interface{}) error {
	// Send a GET request to the URL
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Read the response body
	byteValue, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	// Unmarshal JSON data into the target structure
	err = json.Unmarshal(byteValue, target)
	if err != nil {
		return err
	}

	return nil
}

// GetArtists fetches the list of artists from the API
func GetArtists() ([]models.Artist, error) {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var artists []models.Artist
	err = json.NewDecoder(resp.Body).Decode(&artists)
	if err != nil {
		return nil, err
	}

	return artists, nil
}


// GetLocations fetches location data from the API and unmarshals it into the Location structure
func GetLocations() (models.Location, error) {
	var locationsResponse models.Location
	err := FetchData("https://groupietrackers.herokuapp.com/api/locations", &locationsResponse)
	if err != nil {
		return models.Location{}, err
	}
	return locationsResponse, nil
}

// GetDates fetches date data from the API and unmarshals it into the Date structure
func GetDates() (models.Date, error) {
	var datesResponse models.Date
	err := FetchData("https://groupietrackers.herokuapp.com/api/dates", &datesResponse)
	if err != nil {
		return models.Date{}, err
	}
	return datesResponse, nil
}

// GetRelations fetches relation data from the API and unmarshals it into the Relation structure
func GetRelations() (models.Relation, error) {
	var relationsResponse models.Relation
	err := FetchData("https://groupietrackers.herokuapp.com/api/relation", &relationsResponse)
	if err != nil {
		return models.Relation{}, err
	}
	return relationsResponse, nil
}
