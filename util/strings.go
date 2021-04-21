package util

import "strings"

func SliceContainsString(haystack []string, needle string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}

func JoinStrPointers(sep string, elems ...*string) *string {
	nonNil := []string{}
	for _, elem := range elems {
		if elem != nil {
			nonNil = append(nonNil, *elem)
		}
	}

	// nothing to join
	if len(nonNil) == 0 {
		return nil
	}

	// join & take ptr
	out := strings.Join(nonNil, sep)
	return &out
}
