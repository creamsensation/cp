package cp

import (
	"testing"
	
	"github.com/stretchr/testify/assert"
)

func TestIsFirstCharUpper(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "EmptyString",
			input:    "",
			expected: false,
		},
		{
			name:     "LowerCaseStart",
			input:    "hello",
			expected: false,
		},
		{
			name:     "UpperCaseStart",
			input:    "Hello",
			expected: true,
		},
		{
			name:     "WhiteSpaceStart",
			input:    " hello",
			expected: false,
		},
		{
			name:     "NumberStart",
			input:    "1hello",
			expected: false,
		},
		{
			name:     "SpecialCharStart",
			input:    "!hello",
			expected: false,
		},
	}
	
	for _, testCase := range testCases {
		t.Run(
			testCase.name, func(t *testing.T) {
				result := IsFirstCharUpper(testCase.input)
				assert.Equal(t, testCase.expected, result)
			},
		)
	}
}
