package configs

import (
	"log"
	"os"
	"time"
)

type Config struct {
	HTTP_PORT    string
	DATABASE_URL string
	JWT_SECRET   string
	JWT_TTL      time.Duration
}

func LoadConfig() *Config {
	c := new(Config)
	c.HTTP_PORT = os.Getenv("HTTP_PORT")
	c.DATABASE_URL = os.Getenv("DATABASE_URL")
	c.JWT_SECRET = os.Getenv("JWT_SECRET")
	duration, err := time.ParseDuration(os.Getenv("JWT_TTL"))
	if err != nil {
		log.Println("Invalid JWT TTL")
	}
	c.JWT_TTL = duration
	return c
}
