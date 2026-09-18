package configs

import (
	"log"
	"os"
	"time"
)

type Config struct {
	HTTPPort    string        `env:"HTTP_PORT"`
	DatabaseURL string        `env:"DATABASE_URL"`
	JWTSecret   string        `env:"JWT_SECRET"`
	JWTTTL      time.Duration `env:"JWT_TTL"`
}

func LoadConfig() *Config {
	c := new(Config)
	c.HTTPPort = os.Getenv("HTTP_PORT")
	c.DatabaseURL = os.Getenv("DATABASE_URL")
	c.JWTSecret = os.Getenv("JWT_SECRET")
	duration, err := time.ParseDuration(os.Getenv("JWT_TTL"))
	if err != nil {
		log.Println("Invalid JWT TTL")
	}
	c.JWTTTL = duration
	return c
}
