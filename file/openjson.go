package files

import (
	"encoding/json"
	"io"
	"os"

	"groupie-tracker/models"
)

// openJSON function open a json file from the data folder, and store into the Datas struct
func OpenJSON(file string, data models.Datas) models.Datas {
	jsonfile, _ := os.Open("./data/" + file)
	defer jsonfile.Close()

	// read the data from the json file
	byteValue, _ := io.ReadAll(jsonfile)

	// encoding  of the data into byte
	json.Unmarshal(byteValue, &data)

	return data
}
