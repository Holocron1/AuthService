package configs

import (
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	HTTPPort    string        `env:"HTTP_PORT"`
	DatabaseURL string        `env:"DATABASE_URL"`
	JWTSecret   string        `env:"JWT_SECRET"`
	JWTTTL      time.Duration `env:"JWT_TTL"`
}

func LoadConfig() *Config {
	viper.AutomaticEnv()
	c := new(Config)
	c.HTTPPort = viper.GetString("HTTP_PORT")
	c.DatabaseURL = viper.GetString("DATABASE_URL")
	c.JWTSecret = viper.GetString("JWT_SECRET")
	c.JWTTTL = viper.GetDuration("JWT_TTL")
	return c
}
