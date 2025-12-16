package jwt

import "time"

type Config struct {
	Key            string        `env:"key" envDefault:"secret"`
	ExpireInterval time.Duration `env:"expire_interval" envDefault:"7200s"`
}
