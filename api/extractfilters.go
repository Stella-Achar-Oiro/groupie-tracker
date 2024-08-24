package api

import "net/http"

// ExtractFilters retrieves filter values from the URL query parameters of the HTTP request
func ExtractFilters(req *http.Request) map[string]string {
	filters := make(map[string]string)

	filters["members"] = req.URL.Query().Get("members")
	filters["firstAlbum"] = req.URL.Query().Get("firstAlbum")
	filters["creationDate"] = req.URL.Query().Get("creationDate")
	filters["CitySearch"] = req.URL.Query().Get("CitySearch")

	return filters
}
