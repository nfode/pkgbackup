package main

import (
	"github.com/google/go-cmp/cmp"
	"os"
	"testing"
)

var (
	exampleConfig = Config{
		Hosts: []Host{
			{Name: "host1", File: "host1.txt", IgnoreFile: "host1-ignore.txt", SubscribeTo: []string{"host2"}},
			{Name: "host2", File: "host2.txt", IgnoreFile: "host2-ignore.txt"}},
	}
)

func TestParseConfigFile(t *testing.T) {
	file, _ := os.Open("example/example.yml")
	config := ParseConfigFile(file)
	if !cmp.Equal(exampleConfig, config) {
		t.Fail()
	}
}
