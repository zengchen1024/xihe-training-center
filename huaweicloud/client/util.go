package client

import "strings"

// StrSliceContains checks if a given string is contained in a slice
// When anybody asks why Go needs generics, here you go.
func StrSliceContains(haystack []string, needle string) bool {
	return IsStrContainsSliceElement(needle, haystack, false, true)
}

// IsStrContainsSliceElement returns true if the string exists in given slice or contains in one of slice elements when
// open exact flag. Also you can ignore case for this check.
func IsStrContainsSliceElement(str string, sl []string, ignoreCase, isExcat bool) bool {
	if ignoreCase {
		str = strings.ToLower(str)
	}
	for _, s := range sl {
		if ignoreCase {
			s = strings.ToLower(s)
		}
		if isExcat && s == str {
			return true
		}
		if !isExcat && strings.Contains(str, s) {
			return true
		}
	}
	return false
}
