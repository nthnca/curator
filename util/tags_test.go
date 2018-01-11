package util

import (
	"testing"
)

func TestNormalize(t *testing.T) {
	f := [][][]string{
		{{}, {}},
		{{"1", "2"}, {"1", "2"}},
		{{"1"}, {"1"}},
		{{"1,2"}, {"1", "2"}},
		{{"1,2,3"}, {"1", "2", "3"}},
		{{"1,2,3", "4", "5,6"}, {"1", "2", "3", "4", "5", "6"}},
	}

	for _, x := range f {
		var tags Tags
		tags.A = x[0]
		tags.Normalize()

		if len(tags.A) != len(x[1]) {
			t.Error("Fail: ", 1)
		}

		for i := range tags.A {
			if tags.A[i] != x[1][i] {
				t.Error("Fail: ", 1)
			}
		}
	}
}

func TestValidate(t *testing.T) {
	f := [][][]string{
		{{}, {}, {}},
		{{"1", "2"}, {"1", "2"}, {}},
		{{"1", "2"}, {"1"}, {}},
	}

	for _, x := range f {
		var tags Tags
		tags.A = x[1]
		tags.B = x[2]
		tags.Validate(x[0])
	}
}

func TestModify(t *testing.T) {
	f := [][][]string{
		{{"1"}, {"2"}, {}, {"1"}},
		{{"1"}, {"2"}, {"1", "2"}, {"1"}},
		{{"1"}, {}, {"1", "2"}, {"1", "2"}},
	}

	for _, x := range f {
		var tags Tags
		tags.A = x[0]
		tags.B = x[1]
		out, _ := tags.Modify(x[2])
		exp := x[3]

		if len(out) != len(exp) {
			t.Error("Fail: ", 1)
		}

		for i := range out {
			if out[i] != exp[i] {
				t.Error("Fail: ", 1)
			}
		}
	}
}
