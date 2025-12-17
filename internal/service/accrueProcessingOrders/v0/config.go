package v0

type Config struct {
	Limit int `env:"limit" envDefault:"3"`
}
