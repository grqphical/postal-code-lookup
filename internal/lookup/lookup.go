package lookup

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var provinceAbbreviations map[string]string = map[string]string{
	"British Columbia":             "BC",
	"Alberta":                      "AB",
	"Saskatchewan":                 "SK",
	"Manitoba":                     "MB",
	"Ontario Ottawa":               "ON1",
	"Ontario Greater Toronto Area": "ON2",
	"Ontario Toronto Proper":       "ON3",
	"Ontario South":                "ON4",
	"Ontario North ":               "ON5",
	"Quebec Eastern":               "QC1",
	"Quebec Greater Montreal Area": "QC2",
	"Quebec Eastern Western":       "QC3",
	"New Brunswick":                "NB",
	"Nova Scotia":                  "NS",
	"Prince Edward Island":         "PE",
	"Newfoundland and Labrador":    "NL",
	"Yukon":                        "YT",
	"Northwest Territories":        "NT",
	"Nunavut":                      "NU",
}

type PostalCode struct {
	Urban                      bool   `json:"urban"`
	Province                   string `json:"province"`
	Subdivision                string `json:"subdivision"`
	Municipality               string `json:"municipality"`
	RegionalDistrubutionCentre bool   `json:"regionalDistrubutionCentre"`
	GovernmentBuilding         bool   `json:"governmentBuilding"`
	BusinessReply              bool   `json:"businessReply"`
	PostOffice                 bool   `json:"postOffice"`
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

func stripExtension(filename string) string {
	ext := filepath.Ext(filename)
	return strings.TrimSuffix(filename, ext)
}

func getMunicipalityByFSA(fsa string) (string, error) {
	dataFiles, err := os.ReadDir("data")
	if err != nil {
		return "", errors.New("could not read data directory")
	}

	var provinceFull string
	prov, subdivision, _ := getProvinceSubdivisionFromFSA(fsa)
	if subdivision != "" {
		provinceFull = prov + " " + subdivision
	} else {
		provinceFull = prov
	}

	for _, file := range dataFiles {
		if stripExtension(file.Name()) == provinceAbbreviations[provinceFull] {
			data, err := os.ReadFile(filepath.Join("data", file.Name()))
			if err != nil {
				return "", err
			}

			var jsonData map[string]string
			err = json.Unmarshal(data, &jsonData)
			if err != nil {
				return "", err
			}

			return jsonData[strings.ToUpper(fsa)], nil
		}
	}

	return "", errors.New("this shouldn't happen")
}

func NewPostalCode(postalCode string) (PostalCode, error) {
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
	municipality, err := getMunicipalityByFSA(postalCode[:3])
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
