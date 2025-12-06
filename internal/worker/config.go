package worker

import "time"

type Config struct {
	JobInterval time.Duration `env:"JOB_INTERVAL", envDefault:"3s"`
	JobTimeout  time.Duration `env:"JOB_TIMEOUT", envDefault:"3s"`
}
