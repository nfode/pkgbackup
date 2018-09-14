package main

import (
	"github.com/google/go-cmp/cmp"
	"os"
	"testing"
)

var (
	exampleConfig = Config{
		Hosts: []Host{
			{Name: []string{"host1"}, File: "host1.txt"},
			{Name: []string{"host2"}, File: "host2.txt"}},
	}
)

func TestParseConfigFile(t *testing.T) {
	file, _ := os.Open("example/example.yml")
	config := ParseConfigFile(file)
	if !cmp.Equal(exampleConfig, config) {
		t.Fail()
	}
}
