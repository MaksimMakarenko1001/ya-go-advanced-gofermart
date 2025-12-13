package accrual

import "time"

type Config struct {
	Address        string        `env:"address" envDefault:":8080"`
	Timeout        time.Duration `env:"timeout" envDefault:"3s"`
	ThrottlingRate uint          `env:"throttling_rate" envDefault:"3"`
	MaxRetries     uint16        `env:"max_retries" envDefault:"3"`
	MinDelay       time.Duration `env:"min_delay" envDefault:"1s"`
	DeltaDelay     time.Duration `env:"delta_delay" envDefault:"2s"`
}
