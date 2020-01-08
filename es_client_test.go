package main

import (
	"testing"
)

var cats = []CatEntry{
	CatEntry{
		Index:     "testindex1-2019.07.56",
		StoreSize: "8.5Gb",
	},
}

const (
	ConfigFile = "test_config.json"
)

func TestLoadConfigFile(t *testing.T) {
	c := LoadConfig(ConfigFile)
	if c.Domain != "test-domain" {
		t.Errorf("Domain loaded from config file does not look correct - should be test-domain")
	}
}
