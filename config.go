package main

import "os"

type Config struct {
	Environment   string
	DbTablePrefix string
}

func NewConfig() (*Config, error) {
	c := &Config{
		os.Getenv("ENVIRONMENT"),
		os.Getenv("DB_TABLE_PREFIX"),
	}

	return c, nil
}
