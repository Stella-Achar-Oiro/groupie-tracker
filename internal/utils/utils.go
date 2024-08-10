package utils

import (
	"groupie-tracker/internal/models"
	"sort"
)

func SortArtists(artists []models.Artist, sortBy string) {
	switch sortBy {
	case "name":
		sort.Slice(artists, func(i, j int) bool { return artists[i].Name < artists[j].Name })
	case "date":
		sort.Slice(artists, func(i, j int) bool { return artists[i].CreationDate < artists[j].CreationDate })
	}
}