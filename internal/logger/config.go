package logger

type Config struct {
	Level LogLevel `env:"LEVEL" envDefault:"info"`
}
