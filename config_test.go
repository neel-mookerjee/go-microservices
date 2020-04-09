package main

import (
	"fmt"
	"os"
	"testing"
)

type ConfigTestSuite struct {
	params map[string]string
}

func NewConfigTestSuite() *ConfigTestSuite {
	var params = map[string]string{
		"ENVIRONMENT":     "dev",
		"DB_TABLE_PREFIX": "dev.",
	}

	return &ConfigTestSuite{params: params}
}

func Test_NewConfig(t *testing.T) {
	tc := NewConfigTestSuite()
	tc.loadEnv()
	// Config
	c, _ := NewConfig()

	if c.Environment != tc.params["ENVIRONMENT"] {
		t.Error(fmt.Sprintf("%s expected to be %s but found %s", "c.Environment", tc.params["ENVIRONMENT"], c.Environment))
	}
	if c.DbTablePrefix != tc.params["DB_TABLE_PREFIX"] {
		t.Error(fmt.Sprintf("%s expected to be %s but found %s", "c.DbTablePrefix", tc.params["DB_TABLE_PREFIX"], c.DbTablePrefix))
	}
	tc.removeEnv()
}

func (c *ConfigTestSuite) loadEnv() {
	for k, v := range c.params {
		os.Setenv(k, v)
	}
}

func (c *ConfigTestSuite) removeEnv() {
	for k := range c.params {
		os.Remove(k)
	}
}
