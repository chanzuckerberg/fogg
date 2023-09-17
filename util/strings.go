package util

import (
	"strings"

	"github.com/sirupsen/logrus"
)

func SliceContainsString(haystack []string, needle string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}

func Intersect(slice1, slice2 []string, message *string) (intersection []string) {
	elements := make(map[string]bool, len(slice1))
	for _, item := range slice1 {
		elements[item] = true
	}
	for _, item := range slice2 {
		if elements[item] {
			intersection = append(intersection, item)
		} else {
			if message != nil {
				logrus.Warnf(*message, item)
			}
		}
	}
	return intersection
}

func Union(slice1, slice2 []string) (union []string) {
	elements := make(map[string]bool)
	for _, val := range slice1 {
		elements[val] = true
		union = append(union, val)
	}
	for _, val := range slice2 {
		if !elements[val] {
			union = append(union, val)
		}
	}
	return union
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
