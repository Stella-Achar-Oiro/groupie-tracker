package api

import (
	"fmt"
	"strconv"

	"groupie-tracker/models"
)

func FilterArtists(filters map[string]string, data models.Datas) []models.Artist {
	var filteredArtists []models.Artist

	for artistID, artist := range data.ArtistsData {
		// Check number of members filter
		if membersCount, _ := strconv.Atoi(filters["members"]); len(artist.Members) == membersCount {
			// Check creation date filter
			if creationDate, _ := strconv.Atoi(filters["creationDate"]); artist.CreationDate == creationDate || creationDate == 0 {
				// Check first album filter
				if artist.FirstAlbum == filters["firstAlbum"] || filters["firstAlbum"] == "" {
					// Check city filter
					for _, locationData := range data.LocationsData {
						for _, location := range locationData.Index {
							if location.ID == artistID+1 && (filters["CitySearch"] == "" || contains(location.Locations, filters["CitySearch"])) {
								// Add artist to filtered results and print the artist's name
								filteredArtists = append(filteredArtists, artist)
								fmt.Println(artist.Name)
								break
							}
						}
					}
				}
			}
		}
	}

	return filteredArtists
}

// Helper function to check if a slice contains a specific element
func contains(slice []string, element string) bool {
	for _, value := range slice {
		if value == element {
			return true
		}
	}
	return false
}
