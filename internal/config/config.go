package config

import "os"

type Config struct {
	MongoURI string
}

func Load() (*Config, error) {
	return &Config{
		MongoURI: os.Getenv("MONGO_URI"),
	}, nil
}
