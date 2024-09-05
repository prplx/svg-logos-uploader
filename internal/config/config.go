package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Port              int    `env:"PORT" env-default:"8080"`
	Env               string `env:"ENV" env-default:"production"`
	AdminLogin        string `env:"ADMIN_LOGIN" env-required:"true"`
	AdminPassword     string `env:"ADMIN_PASSWORD" env-required:"true"`
	JWTSecret         string `env:"JWT_SECRET" env-required:"true"`
	GithubAccessToken string `env:"GITHUB_ACCESS_TOKEN" env-required:"true"`
	SessionCookieName string `env:"SESSION_COOKIE_NAME" env-default:"session"`
}

func MustLoadConfig() *Config {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		panic(err)
	}

	return &cfg
}
