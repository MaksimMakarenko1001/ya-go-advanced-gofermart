package v0

type Config struct {
	Limit int `env:"LIMIT" envDefault:"3"`
}
