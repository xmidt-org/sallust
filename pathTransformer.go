// SPDX-FileCopyrightText: 2024 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0
package sallust

import (
	"net/url"
	"os"
)

// PathTransformer is a strategy for altering paths to incorporate
// this package's features.
type PathTransformer struct {
	// Rotation is the optional log rotation configuration.  If supplied,
	// URLs that refer to filesystem paths are altered to be lumberjack URLs.
	Rotation *Rotation

	// Mapping is an optional expansion function passed to os.Expand.  If supplied,
	// this function is used to expand $var and ${var} elements in paths.
	//
	// Any Mapping is always applied to a path first.
	Mapping func(string) string
}

// Transform alters a path to allow for log rotation and expanded variables.
// This method may be passed to ApplyTransform.
func (pt PathTransformer) Transform(path string) (string, error) {
	if pt.Mapping != nil {
		path = os.Expand(path, pt.Mapping)
	}

	if path == "stdout" || path == "stderr" {
		return path, nil
	}

	if pt.Rotation != nil {
		u, err := url.Parse(path)
		if err != nil {
			return path, err
		}

		if len(u.Path) > 0 && (u.Scheme == "" || u.Scheme == "file") {
			path = pt.Rotation.NewURL(u.Path).String()
		}
	}

	return path, nil
}

// ApplyTransform transforms each of a set of paths using the supplied strategy.
// The transformer parameter can be PathTransformer.Transform, or a custom closure.
// This function always returns a newly allocated slice, even if no transformations are done.
// Any error interrupts the transformation, and the transformed slice's contents are undefined.
func ApplyTransform(transformer func(string) (string, error), paths ...string) (transformed []string, err error) {
	if len(paths) == 0 {
		return
	}

	transformed = make([]string, len(paths))
	for i, path := range paths {
		transformed[i], err = transformer(path)
		if err != nil {
			break
		}
	}

	return
}
