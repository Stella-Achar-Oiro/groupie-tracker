package models

type Artist struct {
	ID           int64    `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate uint16   `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`

	LocationsURL    string `json:"locations"`
	ConcertDatesURL string `json:"concertDates"`
	RelationsURL    string `json:"relations"`
}

type Locations struct {
	Index []Location `json:"index"`
}

type Location struct {
	ID        uint64   `json:"id"`
	Locations []string `json:"locations"`
	DatesURL  string   `json:"dates"`
}

type Dates struct {
	Index []Date `json:"index"`
}

type Date struct {
	ID    uint64   `json:"id"`
	Dates []string `json:"dates"`
}

type Relations struct {
	Index []Relation `json:"index"`
}

type Relation struct {
	ID             uint64              `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type ArtistFullData struct {
	Artist         Artist              `json:"artist"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type FullData struct {
	Artists   []Artist  `json:"artists"`
	Locations Locations `json:"locations"`
	Dates     Dates     `json:"dates"`
	Relations Relations `json:"relations"`
}
