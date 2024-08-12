package files

import (
	"fmt"
	"strconv"

	"groupie-tracker/models"
)

func DataToPush(filters map[string]string, data models.Datas) []models.Artists {
	var okFilters models.Datas
	// browses the data from ArtistsData and compares with all the filters
	for id, arti := range data.ArtistsData {
		// try for every artist if it match with filters
		if testMembers, _ := strconv.Atoi(filters["members"]); len(arti.MEMBERS) == testMembers {
			// next filter (Number of members)
			if testDate, _ := strconv.Atoi(filters["creationDate"]); arti.CREA_DATE == testDate || testDate == 0 {
				// next filter (creationDate)
				if arti.FIRST_ALBUM == filters["firstAlbum"] || filters["firstAlbum"] == "" {
					// next filter (FirstAlbum)
					for _, City := range data.LocationsData[id].LOCATIONS {
						// next filter (Locations)
						if filters["CitySearch"] == City || filters["CitySearch"] == "" {
							// Add the result to okFilters of all the filters
							okFilters.ArtistsData = append(okFilters.ArtistsData, data.ArtistsData[id])
							// print if artist matched
							fmt.Println(arti.NAME)
							break
						}
					}
				}
			}
		}
	}

	return okFilters.ArtistsData
}
