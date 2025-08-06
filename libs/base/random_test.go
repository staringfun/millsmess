// Mills Mess
// Licensed under the Mills Mess License Agreement
// See LICENSE.md in the root of this repository.

package base

import (
	"testing"
	"unicode/utf8"
)

func TestGenerateRandomString_Length(t *testing.T) {
	lengths := []int{0, 1, 5, 10, 50, 100}

	for _, l := range lengths {
		result := GenerateRandomString(l)
		if utf8.RuneCountInString(result) != l {
			t.Errorf("Expected string of length %d, got %d", l, utf8.RuneCountInString(result))
		}
	}
}

func TestGenerateRandomString_OnlyUsesAllowedRunes(t *testing.T) {
	allowed := make(map[rune]struct{})
	for _, r := range RandomStringRunes {
		allowed[r] = struct{}{}
	}

	result := GenerateRandomString(1000)
	for _, r := range result {
		if _, ok := allowed[r]; !ok {
			t.Errorf("Character '%c' is not in allowed rune set", r)
		}
	}
}

func TestGenerateRandomString_Uniqueness(t *testing.T) {
	// Generate a few strings and ensure they are (likely) unique.
	// This doesn't guarantee randomness, but detects hardcoded bugs.
	count := 1000
	length := 10
	seen := make(map[string]struct{})

	for i := 0; i < count; i++ {
		s := GenerateRandomString(length)
		if _, exists := seen[s]; exists {
			t.Errorf("Duplicate string generated: %s", s)
		}
		seen[s] = struct{}{}
	}
}
