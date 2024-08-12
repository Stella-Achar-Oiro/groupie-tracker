package models

type Datas struct {
	ArtistsData   []Artists   `json:"Artists_Datas"`
	LocationsData []Locations `json:"Locations_Datas"`
	RelationData  []Relations `json:"Relation_Datas"`
	DatesData     []Dates     `json:"Dates_Datas"`
}