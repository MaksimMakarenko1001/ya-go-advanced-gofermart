package hash

type Config struct {
	Key string `env:"key" envDefault:"secret"`
}
