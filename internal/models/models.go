package models

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

type Datas struct {
    ArtistsData   []Artist   `json:"artists"`
    LocationsData []Location `json:"locations"`
    DatesData     []Date     `json:"dates"`
    RelationsData []Relation `json:"relations"`
}