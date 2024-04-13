package mocks

import (
	"os"

	"go-rest-starter.jtbergman.me/internal/config"
)

// Only exists because flags aren't parsing correctly
func cfg() config.Config {
	cfg := config.Config{}
	cfg.Env = "local"
	cfg.Port = 4000
	cfg.DB.DSN = os.Getenv("TEST_DSN")
	return cfg
}
