package v0

import "time"

type Config struct {
	Timeout time.Duration `env:"timeout" envDefault:"20s"`
}
