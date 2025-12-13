package worker

import "time"

type Config struct {
	JobInterval time.Duration `env:"job_interval" envDefault:"3s"`
	JobTimeout  time.Duration `env:"job_timeout" envDefault:"3s"`
}
