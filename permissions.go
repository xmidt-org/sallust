package sallust

import (
	"fmt"
	"io/fs"
	"strconv"
)

// ParsePermissions parses a nix-style file permissions value.  The value must be a 3-digit
// octal integer with an optional leading zero (0).  The empty string is considered to be 000.
func ParsePermissions(v string) (perms fs.FileMode, err error) {
	switch {
	case len(v) == 0:
		// do nothing.  allow an empty string to map to zero perms

	case len(v) == 3 || (len(v) == 4 && v[0] == '0'):
		var raw uint64
		raw, err = strconv.ParseUint(v, 8, 32)
		if err == nil {
			perms = fs.FileMode(raw)
		} else {
			err = fmt.Errorf("Invalid permissions [%s]: %s", v, err)
		}

	default:
		err = fmt.Errorf("Invalid permissions [%s]: incorrect length", v)
	}

	return
}
