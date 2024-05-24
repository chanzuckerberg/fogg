package util

import (
	"strings"
	"unicode"

	"github.com/sirupsen/logrus"
)

const RootPath = "terraform"

func ConvertToSnake(s string) string {
	var result string
	s = strings.Replace(s, "-", "_", -1)
	v := []rune(s)

	for i := 0; i < len(v); i++ {
		if i != 0 && unicode.IsUpper(v[i]) && (i+1 < len(v) && !unicode.IsUpper(v[i+1])) {
			result += "_"
		}
		result += strings.ToLower(string(v[i]))
	}

	return result
}

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
