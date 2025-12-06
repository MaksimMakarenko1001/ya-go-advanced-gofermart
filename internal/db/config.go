package db

import (
	"errors"
	"strconv"
	"strings"
)

type Config struct {
	DSN        string `json:"-"`
	Host       string `env:"host" envDefault:"localhost"`
	Port       uint16 `env:"port" envDefault:"5432"`
	User       string `env:"user" envDefault:"postgres"`
	Password   string `env:"password" envDefault:"postgres"`
	Name       string `env:"name" envDefault:"postgres"`
	SSLMode    string `env:"ssl_mode" envDefault:"disable"`
	MaxRetries uint16 `env:"max_retries" envDefault:"3"`
}

func (cfg Config) ToDSN() (string, error) {
	if err := cfg.validate(); err != nil {
		return "", err
	}

	kv := []string{
		"host=" + cfg.Host,
		"port=" + strconv.Itoa(int(cfg.Port)),
		"user=" + cfg.User,
		"password=" + cfg.Password,
		"dbname=" + cfg.Name,
		"sslmode=" + cfg.SSLMode,
	}

	return strings.Join(kv, " "), nil

}

func (cfg Config) validate() error {
	var errs []error
	if cfg.Host == "" {
		errs = append(errs, errHostUndefined)
	}
	if cfg.Port == 0 {
		errs = append(errs, errPortUndefined)
	}
	if cfg.Name == "" {
		errs = append(errs, errDBNameUndefined)
	}
	if cfg.User == "" {
		errs = append(errs, errUserUndefined)
	}

	return errors.Join(errs...)
}
