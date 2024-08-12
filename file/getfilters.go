package files

import "net/http"

// getFilters function recovers, from the filters on the MainPage HTML, the various data and stores them into a map (filters)
func GetFilters(r *http.Request) map[string]string {
	filters := make(map[string]string)

	filters["members"] = r.URL.Query().Get("members")
	filters["firstAlbum"] = r.URL.Query().Get("firstAlbum")
	filters["creationDate"] = r.URL.Query().Get("creationDate")
	filters["CitySearch"] = r.URL.Query().Get("CitySearch")

	return filters
}
