package util

import "strings"

func Base(a string) string {
	i := strings.LastIndexAny(a, ".")
	if i >= 0 {
		return a[:i]
	}
	return a
}

func Suffix(a string) string {
	i := strings.LastIndexAny(a, ".")
	if i >= 0 {
		return a[i+1:]
	}
	return ""
}
