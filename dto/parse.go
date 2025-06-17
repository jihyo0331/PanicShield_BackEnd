package dto

import "strconv"

// ParseUint converts a string to uint64, returning an error if invalid.
func ParseUint(s string) (uint64, error) {
	return strconv.ParseUint(s, 10, 64)
}
