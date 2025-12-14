package di

import (
	"flag"
	"log"
	"os"

	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/db"
	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/repository/accrual"
	accrueNewOrders "github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/service/accrueNewOrders/v0"
	getAccrualInfoByOrders "github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/service/getAccrualInfoByOrders/v0"
	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/worker"
	"github.com/caarlos0/env/v6"
)

type diConfig struct {
	AppName    string     `env:"app_name" envDefault:"gofermart"`
	HTTP       HTTPConfig `envPrefix:"HTTP_"`
	DB         db.Config  `envPrefix:"DB_"`
	Repository struct {
		Accrual accrual.Config `envPrefix:"ACCRUAL_"`
	} `envPrefix:"REPOS_"`
	Worker struct {
		AccrueNew        worker.Config `envPrefix:"ACCRUE_NEW_"`
		AccrueProcessing worker.Config `envPrefix:"ACCRUE_PROCESSING_"`
	} `envPrefix:"WORKER_"`
	Service struct {
		AccrueNewOrders        accrueNewOrders.Config        `envPrefix:"ACCRUE_NEW_ORDERS_"`
		GetAccrualInfoByOrders getAccrualInfoByOrders.Config `envPrefix:"GET_ACCRUAL_INFO_BY_ORDERS_"`
	} `envPrefix:"SERVICE_"`
}

func (cfg *diConfig) loadConfig(envPrefix string) {
	cfg.loadDefaults(envPrefix)
	cfg.loadFromArg()
	cfg.loadFromEnv()
}

func (cfg *diConfig) loadDefaults(envPrefix string) {
	if err := env.Parse(cfg, env.Options{Prefix: envPrefix}); err != nil {
		log.Printf("env defaults not ok, %s\n", err.Error())
		return
	}
}

func (cfg *diConfig) loadFromArg() {
	var config diConfig
	flag.StringVar(&config.HTTP.Address, "a", "", "net address")
	flag.StringVar(&config.DB.DSN, "d", "", "data source name")
	flag.StringVar(&config.Repository.Accrual.Address, "r", "", "accrual net address")

	flag.Parse()

	if selfAddr := config.HTTP.Address; selfAddr != "" {
		cfg.HTTP.Address = selfAddr
	}
	if dsn := config.DB.DSN; dsn != "" {
		cfg.DB.DSN = dsn
	}
	if accrualAddr := config.Repository.Accrual.Address; accrualAddr != "" {
		cfg.Repository.Accrual.Address = accrualAddr
	}
}

func (cfg *diConfig) loadFromEnv() {
	if selfAddr := os.Getenv("RUN_ADDRESS"); selfAddr != "" {
		cfg.HTTP.Address = selfAddr
	}
	if dsn := os.Getenv("DATABASE_URI"); dsn != "" {
		cfg.DB.DSN = dsn
	}
	if accrualAddr := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); accrualAddr != "" {
		cfg.Repository.Accrual.Address = accrualAddr
	}

}

type HTTPConfig struct {
	Address string `env:"address" envDefault:":8090"`
}
