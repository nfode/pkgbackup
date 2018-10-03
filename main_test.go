package main

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestGetSystemPackages(t *testing.T) {
	result, err := GetInstalledPackages()
	if err != nil {
		t.Fail()
	}
	if len(result) < 1 {
		t.Fail()
	}
}

func TestReadPackagesFromFile(t *testing.T) {
	packages, err := ReadPackagesFromFile(exampleConfig, "host1", "example")
	if err != nil {
		t.Fail()
	}
	expectedPackages := []string{"firefox-nightly", "spotify"}
	if !cmp.Equal(expectedPackages, packages) {
		fmt.Println("expected: ", expectedPackages)
		fmt.Println("actual: ", packages)
		t.Fail()
	}
}
