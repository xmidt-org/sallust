package sallust

import (
	"encoding/json"
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/natefinch/lumberjack.v2"
)

func testRotationAddQueryValues(t *testing.T) {
	testData := []struct {
		r        Rotation
		expected url.Values
	}{
		{
			r:        Rotation{},
			expected: url.Values{},
		},
		{
			r: Rotation{
				MaxAge:     -1,
				MaxSize:    -1,
				MaxBackups: -1,
				LocalTime:  false,
				Compress:   false,
			},
			expected: url.Values{},
		},
		{
			r: Rotation{
				MaxAge:     156,
				MaxSize:    93723,
				MaxBackups: 483,
				LocalTime:  true,
				Compress:   true,
			},
			expected: url.Values{
				MaxAgeParameter:     []string{"156"},
				MaxSizeParameter:    []string{"93723"},
				MaxBackupsParameter: []string{"483"},
				LocalTimeParameter:  []string{"true"},
				CompressParameter:   []string{"true"},
			},
		},
	}

	for i, record := range testData {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			assert := assert.New(t)
			actual := url.Values{}
			record.r.AddQueryValues(actual)
			assert.Equal(record.expected, actual)
		})
	}
}

func testRotationNewURL(t *testing.T) {
	testData := []struct {
		r        Rotation
		path     string
		expected string
	}{
		{
			r:        Rotation{},
			path:     "/var/log/foo.json",
			expected: "lumberjack:///var/log/foo.json",
		},
		{
			r: Rotation{
				MaxAge:     -1,
				MaxSize:    -1,
				MaxBackups: -1,
				LocalTime:  false,
				Compress:   false,
			},
			path:     "/defaults.log",
			expected: "lumberjack:///defaults.log",
		},
		{
			r: Rotation{
				MaxAge:     77,
				MaxSize:    459,
				MaxBackups: 1774,
				LocalTime:  true,
				Compress:   true,
			},
			path:     "/test.json",
			expected: "lumberjack:///test.json?compress=true&localTime=true&maxAge=77&maxBackups=1774&maxSize=459",
		},
	}

	for i, record := range testData {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var (
				assert  = assert.New(t)
				require = require.New(t)
				actual  = record.r.NewURL(record.path)
			)

			require.NotNil(actual)

			// url.Values.Encode() produces a query string sorted by key,
			// so comparisons to expected values should be consistent.
			assert.Equal(record.expected, actual.String())
		})
	}
}

func TestRotation(t *testing.T) {
	t.Run("AddQueryValues", testRotationAddQueryValues)
	t.Run("NewURL", testRotationNewURL)
}

func testLumberjackSync(t *testing.T) {
	assert := assert.New(t)
	assert.NotPanics(func() {
		// Sync should be a nop, and shouldn't invoke anything on the *lumberjack.Logger
		assert.NoError(Lumberjack{}.Sync())
	})
}

func TestLumberjack(t *testing.T) {
	t.Run("Sync", testLumberjackSync)
}

func testNewLumberjackSinkSuccess(t *testing.T) {
	testData := []struct {
		actualURL string
		expected  *lumberjack.Logger
	}{
		{
			actualURL: "lumberjack:///var/log/foo",
			expected: &lumberjack.Logger{
				Filename: "/var/log/foo",
			},
		},
		{
			actualURL: "lumberjack:///log.json?maxAge=12&maxBackups=72&maxSize=19000&compress=true&localTime=true",
			expected: &lumberjack.Logger{
				Filename:   "/log.json",
				MaxAge:     12,
				MaxBackups: 72,
				MaxSize:    19000,
				Compress:   true,
				LocalTime:  true,
			},
		},
	}

	for i, record := range testData {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var (
				assert  = assert.New(t)
				require = require.New(t)
			)

			u, err := url.Parse(record.actualURL)
			require.NoError(err)
			require.NotNil(u)

			s, err := NewLumberjackSink(u)
			require.NoError(err)
			require.NotNil(s)

			actual, ok := s.(Lumberjack)
			require.True(ok)

			// we only want to compare the public fields, so use JSON marshalling as a filter
			expectedJSON, err := json.Marshal(record.expected)
			require.NoError(err)

			actualJSON, err := json.Marshal(actual.Logger)
			require.NoError(err)

			assert.JSONEq(string(expectedJSON), string(actualJSON))
		})
	}
}

func testNewLumberjackSinkInvalidURL(t *testing.T) {
	testData := []url.URL{
		{
			Path:     "/test",
			RawQuery: "t=%^X",
		},
		{
			Path:     "/test",
			RawQuery: "maxAge=thisisnotavalidint",
		},
		{
			Path:     "/test",
			RawQuery: "maxBackups=thisisnotavalidint",
		},
		{
			Path:     "/test",
			RawQuery: "maxSize=thisisnotavalidint",
		},
		{
			Path:     "/test",
			RawQuery: "localTime=thisisnotavalidbool",
		},
		{
			Path:     "/test",
			RawQuery: "compress=thisisnotavalidbool",
		},
	}

	for i, badURL := range testData {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var (
				assert = assert.New(t)
				s, err = NewLumberjackSink(&badURL)
			)

			assert.Error(err)
			assert.Nil(s)
		})
	}
}

func TestNewLumberjackSink(t *testing.T) {
	t.Run("Success", testNewLumberjackSinkSuccess)
	t.Run("InvalidURL", testNewLumberjackSinkInvalidURL)
}
