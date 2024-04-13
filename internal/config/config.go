package config

import (
	"flag"
	"fmt"
	"log"
	"os"
)

// ============================================================================
// Constants
// ============================================================================

const (
	EnvLocal = "local"
	EnvDev   = "dev"
	EnvProd  = "prod"
)

// ============================================================================
// Config
// ============================================================================

// Defines the configuration provided via command line args
type Config struct {
	Env  string
	Port int
	DB   struct {
		DSN string
	}
	SMTP struct {
		Host     string
		Port     int
		Username string
		Password string
		Sender   string
	}
}

// Create validated config
func New() Config {
	cfg := &Config{}

	// Env
	flag.StringVar(&cfg.Env, "env", "", "Development environment")

	// Server
	flag.IntVar(&cfg.Port, "port", 0, "API server port")

	// Database
	flag.StringVar(&cfg.DB.DSN, "db-dsn", "", "Postgres DSN")

	// SMTP
	flag.StringVar(&cfg.SMTP.Host, "smtp-host", "", "SMTP host")
	flag.IntVar(&cfg.SMTP.Port, "smtp-port", 0, "SMTP port")
	flag.StringVar(&cfg.SMTP.Username, "smtp-username", "", "SMTP username")
	flag.StringVar(&cfg.SMTP.Password, "smtp-password", "", "SMTP password")
	flag.StringVar(&cfg.SMTP.Sender, "smtp-sender", "", "SMTP sender")

	// Version
	displayVersion := flag.Bool("version", false, "Display version and exit")

	flag.Parse()

	if *displayVersion {
		fmt.Printf("Version:\t%s\n", version())
		os.Exit(0)
	}

	// Validate
	if ok, err := cfg.validate(); !ok {
		log.Fatalln(err)
	}

	return *cfg
}

// Returns true if the local environment is used
func (config *Config) IsLocal() bool {
	return config.Env == EnvLocal
}

// ============================================================================
// Validation
// ============================================================================

// Validate config and return ok or error
func (config *Config) validate() (bool, string) {
	// Validate env
	switch config.Env {
	case EnvLocal, EnvDev, EnvProd:
		break

	default:
		return false, fmt.Sprintf("Missing env flag (%s | %s | %s)", EnvLocal, EnvDev, EnvProd)
	}

	// Validate ints
	switch 0 {
	case config.Port:
		return false, "Missing port flag"

	case config.SMTP.Port:
		if !config.IsLocal() {
			return false, "Missing smtp-port flag"
		}
	}

	// Validate strings
	if !config.IsLocal() {
		switch "" {
		case config.DB.DSN:
			return false, "Missing db-dsn flag"

		case config.SMTP.Host:
			return false, "Missing smtp-host flag"

		case config.SMTP.Username:
			return false, "Missing smtp-username flag"

		case config.SMTP.Password:
			return false, "Missing smtp-password flag"

		case config.SMTP.Sender:
			return false, "Missing smtp-sender flag"
		}
	}

	return true, ""
}
