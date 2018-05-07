package util

import (
	"log"
	"strings"
)

type Tags struct {
	A []string
	B []string
}

func (tags *Tags) Normalize() {
	normalize := func(arr []string) []string {
		var tmp []string
		for _, t := range arr {
			tmp = append(tmp, strings.Split(t, ",")...)
		}
		return tmp
	}

	tags.A = normalize(tags.A)
	tags.B = normalize(tags.B)
}

func (tags *Tags) Validate(allowed []string) {
	allowedx := make(map[string]bool)

	for _, t := range allowed {
		allowedx[t] = true
	}

	valid := func(arr []string) {
		for _, x := range arr {
			if !allowedx[x] {
				log.Fatalf("Invalid tag: %s", x)
			}
		}
	}

	valid(tags.A)
	valid(tags.B)

	duplicates := map[string]bool{}

	dups := func(arr []string) {
		for _, x := range arr {
			if duplicates[x] {
				log.Fatalf("Duplicate tag: %s", x)
			}
			duplicates[x] = true
		}
	}

	dups(tags.A)
	dups(tags.B)
}

func (tags *Tags) Modify(labelList []string) ([]string, bool) {
	changed := false
	tagMap := make(map[string]bool)

	for _, t := range labelList {
		tagMap[t] = true
	}
	for _, t := range tags.A {
		_, has := tagMap[t]
		changed = changed || !has
		tagMap[t] = true
	}
	for _, t := range tags.B {
		_, has := tagMap[t]
		changed = changed || has
		delete(tagMap, t)
	}
	var rv []string
	for t, _ := range tagMap {
		rv = append(rv, t)
	}
	return rv, changed
}

func (tags *Tags) Match(labelList []string) bool {
	for _, a := range tags.A {
		match := false
		for _, t := range labelList {
			if a == t {
				match = true
				break
			}
		}
		if !match {
			return false
		}
	}
	for _, b := range tags.B {
		for _, t := range labelList {
			if b == t {
				return false
			}
		}
	}
	return true
}

/*
func deleted(lbls []string) bool {
	for _, x := range lbls {
		if x == "keep" {
			return false
		}
	}
	return true
}
*/
