package v0

import "time"

type Config struct {
	LockInterval time.Duration `env:"lock_interval" envDefault:"3s"`
}
