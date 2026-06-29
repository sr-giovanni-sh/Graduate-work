package config

import "os"

type AppConfig struct {
	Password string
	Port     string
	DBFile   string
}

func Load() AppConfig {
	return AppConfig{
		Password: os.Getenv("TODO_PASSWORD"),
		Port:     os.Getenv("TODO_PORT"),
		DBFile:   os.Getenv("TODO_DBFILE"),
	}
}
