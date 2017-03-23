package util

import (
	"log"
	"reflect"
	"testing"
)

func TestGetRandomSet(t *testing.T) {
	if v := GetRandomSet(4, 4); !reflect.DeepEqual(v, []int{0, 1, 2, 3}) {
		t.Error("Fail: ", v)
	}

	/*
		if v := GetRandomSet(4, 5); err == nil || v != nil {
			t.Error("Fail: ", v)
		}
	*/

	v := GetRandomSet(400, 5)
	log.Printf("GetRandomSet(400, 5): %v", v)
}
