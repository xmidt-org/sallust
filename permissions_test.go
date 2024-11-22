package sallust

import (
	"io/fs"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ParsePermissionsTestSuite struct {
	suite.Suite
}

func (suite *ParsePermissionsTestSuite) TestValid() {
	testCases := []struct {
		value    string
		expected fs.FileMode
	}{
		{value: "", expected: 0},
		{value: "644", expected: 0644},
		{value: "0644", expected: 0644},
		{value: "666", expected: 0666},
		{value: "0666", expected: 0666},
		{value: "600", expected: 0600},
		{value: "0600", expected: 0600},
		{value: "000", expected: 0000},
		{value: "0000", expected: 0000},
	}

	for _, testCase := range testCases {
		suite.Run(testCase.value, func() {
			actual, err := ParsePermissions(testCase.value)
			suite.Require().NoError(err)
			suite.Equal(testCase.expected, actual)
		})
	}
}

func (suite *ParsePermissionsTestSuite) TestInvalid() {
	testCases := []string{
		"0",
		"463213",
		"900",
		"0900",
		"090",
		"0090",
		"009",
		"0009",
		"x000",
		"this is definitely invalid",
	}

	for _, testCase := range testCases {
		suite.Run(testCase, func() {
			_, err := ParsePermissions(testCase)
			suite.Error(err)
		})
	}
}

func TestParsePermissions(t *testing.T) {
	suite.Run(t, new(ParsePermissionsTestSuite))
}
