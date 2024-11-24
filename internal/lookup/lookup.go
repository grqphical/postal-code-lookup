package lookup

import (
	"database/sql"
	"errors"
	"regexp"
	"strings"
)

// PostalCode model info
// @Description Contains the extracted data from the given postal code
type PostalCode struct {
	Urban                      bool   `json:"urban"`                      // whether or not the postal code is in an urban center
	Province                   string `json:"province"`                   // what province/territory the postal code is in
	Subdivision                string `json:"subdivision"`                // for the bigger provinces, which area of the province is it in
	Municipality               string `json:"municipality"`               // what municipality the postal code is located in
	RegionalDistrubutionCentre bool   `json:"regionalDistrubutionCentre"` // is the postal code a Canada post distribution centre
	GovernmentBuilding         bool   `json:"governmentBuilding"`         // does the postal code belong to a building owned by the Canadian government
	BusinessReply              bool   `json:"businessReply"`              // is this code used as part of the Business Reply system
	PostOffice                 bool   `json:"postOffice"`                 // is this code a Canada Post post office
}

// IsValidPostalCode checks if a given string is a valid Canadian postal code
func IsValidPostalCode(postalCode string) bool {
	pattern := `^[abceghjklmnprstvxy]\d[abcdefghjklmnprstvxyz][ ]?\d[abceghjklmnprstvxywz]\d`

	re := regexp.MustCompile(pattern)

	return re.MatchString(postalCode)
}

// Gets the associated province and subdivision (if applicable) based on the given Forward Sortition Area (first three characters of code)
// If the FSA is invalid, it returns an error
func getProvinceSubdivisionFromFSA(fsa string) (string, string, error) {
	switch fsa[0] {
	case 'a':
		return "Newfoundland", "", nil
	case 'b':
		return "Nova Scotia", "", nil
	case 'c':
		return "Prince Edward Island", "", nil
	case 'e':
		return "New Brunswick", "", nil
	case 'g':
		return "Quebec", "Eastern", nil
	case 'h':
		return "Quebec", "Greater Montreal Area", nil
	case 'j':
		return "Quebec", "Western", nil
	case 'k':
		return "Ontario", "Ottawa", nil
	case 'l':
		return "Ontario", "Greater Toronto Area", nil
	case 'm':
		return "Ontario", "Toronto Proper", nil
	case 'n':
		return "Ontario", "South", nil
	case 'p':
		return "Ontarto", "North", nil
	case 'r':
		return "Manitoba", "", nil
	case 's':
		return "Saskatchewan", "", nil
	case 't':
		return "Alberta", "", nil
	case 'v':
		return "British Columbia", "", nil
	case 'x':
		if fsa == "x0a" || fsa == "x0b" || fsa == "x0c" {
			return "Nunavut", "", nil
		}

		return "North West Territories", "", nil
	case 'y':
		return "Yukon", "", nil
	default:
		return "", "", errors.New("invalid fsa")
	}
}

// Gets the municipality for the given Forward Sorition Area from the database
func getMunicipalityByFSA(fsa string, db *sql.DB) (string, error) {
	var municipality string
	row := db.QueryRow("SELECT municipality FROM Municipalities WHERE fsa = ?", strings.ToUpper(fsa))
	err := row.Scan(&municipality)

	return municipality, err
}

// Validates a postal code then extracts all relevant information from the code. After validation and extraction,
// this function returns a struct containing all the information extracted.
func NewPostalCode(postalCode string, db *sql.DB) (PostalCode, error) {
	var postalCodeObj PostalCode
	postalCode = strings.ToLower(postalCode)
	postalCode = strings.ReplaceAll(postalCode, " ", "")

	if !IsValidPostalCode(postalCode) {
		return postalCodeObj, errors.New("invalid postal code")
	}

	province, subdivision, err := getProvinceSubdivisionFromFSA(postalCode[:3])
	if err != nil {
		return postalCodeObj, err
	}

	postalCodeObj.Province = province
	postalCodeObj.Subdivision = subdivision
	municipality, err := getMunicipalityByFSA(postalCode[:3], db)
	if err != nil {
		return postalCodeObj, err
	}
	postalCodeObj.Municipality = municipality

	if postalCode[1] != '0' {
		postalCodeObj.Urban = true
	}
	if postalCode[3:6] == "9z9" {
		postalCodeObj.BusinessReply = true
	} else if postalCode[3:6] == "9z0" {
		postalCodeObj.RegionalDistrubutionCentre = true
	}

	if postalCode[:3] == "k1a" {
		postalCodeObj.GovernmentBuilding = true
	}

	if postalCode[5] == '0' {
		postalCodeObj.PostOffice = true
	}

	return postalCodeObj, nil
}
