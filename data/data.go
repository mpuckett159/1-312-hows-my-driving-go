package data

import (
	"1-312-hows-my-driving-go/csvmap"
)

// BadgeData : map[string]string of officer badge data
var BadgeData, _ = csvmap.CSVFileToMap("data/spd-badges.csv")

// SearchCSVMap : Function to search CSV data stored as []map[string]string
func SearchCSVMap(searchMap map[string]string) (returnMaps []map[string]string) {
	for i := range BadgeData {
		for key, value := range searchMap {
			if BadgeData[i][key] == value {
				returnMaps = append(returnMaps, BadgeData[i])
			}
		}
	}
	return
}
