package main

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestComparePackages(t *testing.T) {
	first := []string{"1", "2", "3"}
	second := []string{"1", "2", "4"}
	expected := CompareResult{
		Added:     []string{"3"},
		Unchanged: []string{"1", "2"},
		Removed:   []string{"4"},
	}
	result := ComparePackages(first, second)
	if !cmp.Equal(expected, result) {
		fmt.Println("expected: ", expected)
		fmt.Println("actual: ", result)
		t.Fail()
	}
}
