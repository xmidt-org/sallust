// SPDX-FileCopyrightText: 2024 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0
package sallust

import (
	"errors"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testPathTransformerSuccess(t *testing.T) {
	const testValue = "test"

	testData := []struct {
		pt       PathTransformer
		path     string
		expected string
	}{
		{
			pt:       PathTransformer{},
			path:     "",
			expected: "",
		},
		{
			pt:       PathTransformer{},
			path:     "stdout",
			expected: "stdout",
		},
		{
			pt:       PathTransformer{},
			path:     "stderr",
			expected: "stderr",
		},
		{
			pt: PathTransformer{
				Mapping: func(v string) string {
					if v == testValue {
						return "stdout"
					}

					return ""
				},
			},
			path:     "${test}",
			expected: "stdout",
		},
		{
			pt:       PathTransformer{},
			path:     "/var/log/app.json",
			expected: "/var/log/app.json",
		},
		{
			pt: PathTransformer{
				Rotation: &Rotation{
					MaxAge:   8,
					MaxSize:  45,
					Compress: true,
				},
			},
			path:     "/var/log/log.json",
			expected: "lumberjack:///var/log/log.json?compress=true&maxAge=8&maxSize=45",
		},
		{
			pt: PathTransformer{
				Rotation: &Rotation{
					MaxSize:    47,
					MaxBackups: 5,
				},
			},
			path:     "file:///var/log/log.json",
			expected: "lumberjack:///var/log/log.json?maxBackups=5&maxSize=47",
		},
		{
			pt: PathTransformer{
				Rotation: &Rotation{
					MaxAge:  10,
					MaxSize: 150,
				},
				Mapping: func(v string) string {
					if v == testValue {
						return "/var/log"
					}

					return ""
				},
			},
			path:     "$test/log.json",
			expected: "lumberjack:///var/log/log.json?maxAge=10&maxSize=150",
		},
		{
			pt: PathTransformer{
				Rotation: &Rotation{
					MaxAge:     417,
					MaxBackups: 3,
					LocalTime:  true,
				},
				Mapping: func(v string) string {
					if v == testValue {
						return "/var/log"
					}

					return ""
				},
			},
			path:     "file://$test/log.json",
			expected: "lumberjack:///var/log/log.json?localTime=true&maxAge=417&maxBackups=3",
		},
	}

	for i, record := range testData {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var (
				assert      = assert.New(t)
				actual, err = record.pt.Transform(record.path)
			)

			assert.NoError(err)
			assert.Equal(record.expected, actual)
		})
	}
}

func testPathTransformerInvalidURL(t *testing.T) {
	var (
		assert = assert.New(t)
		pt     = PathTransformer{
			Rotation: &Rotation{
				MaxSize: 45,
			},
		}
	)

	_, err := pt.Transform("^*@&#(&%@#%XX")
	assert.Error(err)
}

func TestPathTransformer(t *testing.T) {
	t.Run("Success", testPathTransformerSuccess)
	t.Run("InvalidURL", testPathTransformerInvalidURL)
}

func testApplyTransformSuccess(t *testing.T) {
	var (
		assert = assert.New(t)

		transformer = func(v string) (string, error) {
			return "transformed:" + v, nil
		}
	)

	result, err := ApplyTransform(transformer)
	assert.Empty(result)
	assert.NoError(err)

	result, err = ApplyTransform(transformer, []string{}...)
	assert.Empty(result)
	assert.NoError(err)

	result, err = ApplyTransform(transformer, "1")
	assert.Equal(result, []string{"transformed:1"})
	assert.NoError(err)

	result, err = ApplyTransform(transformer, "1", "2")
	assert.Equal(result, []string{"transformed:1", "transformed:2"})
	assert.NoError(err)

	result, err = ApplyTransform(transformer, "1", "2", "3", "4", "5")
	assert.Equal(
		result,
		[]string{"transformed:1", "transformed:2", "transformed:3", "transformed:4", "transformed:5"},
	)

	assert.NoError(err)

	result, err = ApplyTransform(
		func(string) (string, error) { return "", errors.New("should not be called") },
	)

	assert.Empty(result)
	assert.NoError(err)
}

func testApplyTransformFailure(t *testing.T) {
	testData := []struct {
		transformer func(string) (string, error)
		paths       []string
	}{
		{
			transformer: func(string) (string, error) {
				return "", errors.New("expected")
			},
			paths: []string{"1"},
		},
		{
			transformer: func(v string) (string, error) {
				if v == "1" {
					return "transformed", nil
				}

				return "", errors.New("expected")
			},
			paths: []string{"1", "2"},
		},
		{
			transformer: func(v string) (string, error) {
				if v != "3" {
					return "transformed", nil
				}

				return "", errors.New("expected")
			},
			paths: []string{"1", "2", "3"},
		},
	}

	for i, record := range testData {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			assert := assert.New(t)
			_, err := ApplyTransform(record.transformer, record.paths...)
			assert.Error(err)
		})
	}
}

func TestApplyTransform(t *testing.T) {
	t.Run("Success", testApplyTransformSuccess)
	t.Run("Failure", testApplyTransformFailure)
}
