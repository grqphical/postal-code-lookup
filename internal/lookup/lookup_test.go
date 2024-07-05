package lookup

import (
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
		assert.Equal(t, testCase.expected, isValidPostalCode(strings.ToLower(testCase.postalCode)), "Failed on: %s", testCase.postalCode)
	}

}

func TestNewPostalCode(t *testing.T) {
	testCases := []struct {
		postalCode     string
		expected       PostalCode
		expectingError bool
		expectedErrMsg string
	}{
		{"k1a 0b1", PostalCode{Province: "Ontario", Subdivision: "Ottawa", GovernmentBuilding: true, Urban: true}, false, ""},
		{"b2c 3e4", PostalCode{Province: "Nova Scotia", Urban: true}, false, ""},
		{"g1x 2z0", PostalCode{Province: "Quebec", Subdivision: "Eastern", Urban: true, PostOffice: true}, false, ""},
		{"k1a1b2", PostalCode{Province: "Ontario", Subdivision: "Ottawa", GovernmentBuilding: true, Urban: true}, false, ""},
		{"b2c9z9", PostalCode{Province: "Nova Scotia", Urban: true, BusinessReply: true}, false, ""},
		{"x0a1b2", PostalCode{Province: "Nunavut"}, false, ""},
		{"k1a-0b1", PostalCode{}, true, "invalid postal code"},
		{"invalid", PostalCode{}, true, "invalid postal code"},
	}

	for _, testCase := range testCases {
		postalCodeObj, err := NewPostalCode(testCase.postalCode)

		if testCase.expectingError {
			assert.NotNil(t, err)
			assert.Equal(t, testCase.expectedErrMsg, err.Error())
		} else {
			assert.Nil(t, err)
			assert.Equal(t, testCase.expected, postalCodeObj, testCase.postalCode)
		}
	}
}
