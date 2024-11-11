package lookup

import (
	"database/sql"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidPostalCode(t *testing.T) {
	testCases := []struct {
		postalCode string
		expected   bool
	}{
		{"K1A 0B1", true},  // Valid
		{"Z9Z 9Z9", false}, // Invalid (Z as first character)
		{"123 456", false}, // Invalid (only digits)
		{"A1A1A1", true},   // Valid (no space)
		{"K1A-0B1", false}, // Invalid (contains hyphen)
		{"", false},        // Invalid (empty string)
		{"K1A0B1", true},   // Valid (no space)
		{"A1A-1A1", false}, // Invalid (contains hyphen)
		{"A1A A1A", false}, // Invalid (wrong format)
		{"X0X 0X0", true},  // Valid (edge case with X)
		{"A1 1A1", false},  // Invalid (missing character),
		{"foobaff", false}, // Invalid (not a postal code at all)
	}

	for _, testCase := range testCases {
		assert.Equal(t, testCase.expected, IsValidPostalCode(strings.ToLower(testCase.postalCode)), "Failed on: %s", testCase.postalCode)
	}

}

func TestNewPostalCode(t *testing.T) {
	testCases := []struct {
		postalCode     string
		expected       PostalCode
		expectingError bool
		expectedErrMsg string
	}{
		{"k1a 0b1", PostalCode{Province: "Ontario", Subdivision: "Ottawa", GovernmentBuilding: true, Urban: true, Municipality: "Government of Canada Ottawa and Gatineau offices (partly in QC)"}, false, ""},
		{"b2c 3e4", PostalCode{Province: "Nova Scotia", Urban: true, Municipality: "Iona"}, false, ""},
		{"g1x 2z0", PostalCode{Province: "Quebec", Subdivision: "Eastern", Urban: true, PostOffice: true, Municipality: "Quebec City West Sainte-Foy"}, false, ""},
		{"k1a1b2", PostalCode{Province: "Ontario", Subdivision: "Ottawa", GovernmentBuilding: true, Urban: true, Municipality: "Government of Canada Ottawa and Gatineau offices (partly in QC)"}, false, ""},
		{"b2c9z9", PostalCode{Province: "Nova Scotia", Urban: true, BusinessReply: true, Municipality: "Iona"}, false, ""},
		{"x0a1b2", PostalCode{Province: "Nunavut", Municipality: "Qikiqtaaluk Region 0A0: Arctic Bay 0B0: Qikiqtarjuaq 0C0: Kinngait 0E0: Clyde River 0G0: Eureka 0H0: Iqaluit 0J0: Grise Fiord 0K0: Sanirajak 0L0: Igloolik 0N0: Kimmirut 0R0: Pangnirtung 0S0: Pond Inlet 0V0: Resolute 0W0: Sanikiluaq 1H0: Iqaluit 2H0: Iqaluit 3H0: Iqaluit"}, false, ""},
		{"k1a-0b1", PostalCode{}, true, "invalid postal code"},
		{"invalid", PostalCode{}, true, "invalid postal code"},
	}

	conn, err := sql.Open("sqlite3", os.Getenv("DATABASE_URL"))
	assert.Nil(t, err)

	for _, testCase := range testCases {
		postalCodeObj, err := NewPostalCode(testCase.postalCode, conn)

		if testCase.expectingError {
			assert.NotNil(t, err, testCase.postalCode)
			assert.Equal(t, testCase.expectedErrMsg, err.Error())
		} else {
			assert.Nil(t, err, testCase.postalCode)
			assert.Equal(t, testCase.expected, postalCodeObj, testCase.postalCode)
		}
	}
}
