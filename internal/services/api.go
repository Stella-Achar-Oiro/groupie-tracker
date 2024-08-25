package services

import (
    "encoding/json"
    "net/http"
    "groupie-tracker/internal/models"
)

const (
    ArtistsAPI   = "https://groupietrackers.herokuapp.com/api/artists"
    LocationsAPI = "https://groupietrackers.herokuapp.com/api/locations"
    DatesAPI     = "https://groupietrackers.herokuapp.com/api/dates"
    RelationsAPI = "https://groupietrackers.herokuapp.com/api/relation"
)

func FetchData() (*models.Datas, error) {
    artists, err := fetchArtists()
    if err != nil {
        return nil, err
    }

    locations, err := fetchLocations()
    if err != nil {
        return nil, err
    }

    dates, err := fetchDates()
    if err != nil {
        return nil, err
    }

    relations, err := fetchRelations()
    if err != nil {
        return nil, err
    }

    return &models.Datas{
        ArtistsData:   artists,
        LocationsData: []models.Location{locations},
        DatesData:     []models.Date{dates},
        RelationsData: []models.Relation{relations},
    }, nil
}

func fetchArtists() ([]models.Artist, error) {
    resp, err := http.Get(ArtistsAPI)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var artists []models.Artist
    err = json.NewDecoder(resp.Body).Decode(&artists)
    return artists, err
}

func fetchLocations() (models.Location, error) {
    resp, err := http.Get(LocationsAPI)
    if err != nil {
        return models.Location{}, err
    }
    defer resp.Body.Close()

    var locations models.Location
    err = json.NewDecoder(resp.Body).Decode(&locations)
    return locations, err
}

func fetchDates() (models.Date, error) {
    resp, err := http.Get(DatesAPI)
    if err != nil {
        return models.Date{}, err
    }
    defer resp.Body.Close()

    var dates models.Date
    err = json.NewDecoder(resp.Body).Decode(&dates)
    return dates, err
}

func fetchRelations() (models.Relation, error) {
    resp, err := http.Get(RelationsAPI)
    if err != nil {
        return models.Relation{}, err
    }
    defer resp.Body.Close()

    var relations models.Relation
    err = json.NewDecoder(resp.Body).Decode(&relations)
    return relations, err
}