package sallust

import (
	"errors"
	"io/fs"
)

var (
	// ErrInvalidPermissions is returned by ParsePermissions to indicate a bad permissions value.
	ErrInvalidPermissions = errors.New("Invalid permissions")
)

func accumulate(v byte, factor int, perms *fs.FileMode) (ok bool) {
	if ok = v >= '0' && v <= '7'; ok {
		*perms += (fs.FileMode(int(v-'0') * factor))
	}

	return
}

// ParsePermissions parses a nix-style file permissions value.  The value must be a 3-digit
// octal integer with an optional leading zero (0).  The empty string is considered to be 000.
func ParsePermissions(v string) (perms fs.FileMode, err error) {
	switch {
	// allows for an unset configuration value, which means just take the underlying defaults
	case len(v) == 0:
		return

	case len(v) < 3:
		fallthrough

	case len(v) > 4:
		fallthrough

	// if 4 characters, the first character must be a zero (0)
	case len(v) == 4 && v[0] != '0':
		err = ErrInvalidPermissions

	case !accumulate(v[len(v)-1], 1, &perms):
		err = ErrInvalidPermissions

	case !accumulate(v[len(v)-2], 8, &perms):
		err = ErrInvalidPermissions

	case !accumulate(v[len(v)-3], 64, &perms):
		err = ErrInvalidPermissions
	}

	return
}
