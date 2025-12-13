package di

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/db"
	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/repository/accrual"
	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/repository/lock"
	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/repository/pg"
	accrueNewOrders "github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/service/accrueNewOrders/v0"
	getAccrualInfoByOrders "github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/service/getAccrualInfoByOrders/v0"
	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/worker"
	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/pkg/backoff"
)

type DI struct {
	config       *diConfig
	repositories struct {
		accrual *accrual.Repository
		pg      *pg.Repository
		lock    *lock.Repository
	}
	services struct {
		included struct {
			getAccrualInfoByOrders *getAccrualInfoByOrders.Service
		}
		accrueNewOrders *accrueNewOrders.Service
	}
	workers struct {
		accrueNew        *worker.Worker
		accrueProcessing *worker.Worker
	}
	// api struct {
	// 	external *handler.API
	// }
	infr struct {
		db *db.PGConnect
	}
}

func (di *DI) Init(envPrefix string) {
	di.config = &diConfig{}
	di.config.loadConfig(envPrefix)

	di.initDB()
	di.initRepositories()
	di.initServices()
	di.initWorkers()
	// di.initAPI()
}

func (di *DI) initDB() {
	var err error
	di.infr.db, err = db.New(
		di.config.DB,
		backoff.NewLinearBackoff(
			backoff.NewBackoff(db.ClassifyPgError, di.config.DB.MaxRetries, di.config.DB.MinDelay),
			di.config.DB.DeltaDelay,
		),
	)
	if err != nil {
		log.Println("db init not ok,", err.Error())
	}
}

func (di *DI) initRepositories() {
	di.repositories.accrual = accrual.New(di.config.Repository.Accrual, backoff.NewLinearBackoff(
		backoff.NewBackoff(accrual.ClassifyHTTPError, di.config.Repository.Accrual.MaxRetries, di.config.DB.MinDelay),
		di.config.DB.DeltaDelay,
	))
	di.repositories.lock = lock.New(di.infr.db)
	di.repositories.pg = pg.New(di.infr.db)
}

func (di *DI) initServices() {
	di.services.included.getAccrualInfoByOrders = getAccrualInfoByOrders.New(di.config.Service.GetAccrualInfoByOrders, di.repositories.accrual)

	di.services.accrueNewOrders = accrueNewOrders.New(di.config.Service.AccrueNewOrders, di.repositories.pg, di.services.included.getAccrualInfoByOrders)
}

func (di *DI) initWorkers() {
	di.workers.accrueNew = worker.New(
		di.config.Worker.AccrueNew,
		di.repositories.lock,
		di.config.AppName,
		"accrue_new",
		"",
		di.services.accrueNewOrders.Do,
	)
}

// func (di *DI) initAPI() {
// 	di.api.external = handler.New(
// 		logger.New(di.config.Logger),
// 		di.services.updateFlatService,
// 		di.services.updateBatchService,
// 		di.services.updateService,
// 		di.services.getFlatService,
// 		di.services.getService,
// 		di.services.listMetricService,
// 		di.services.dumpMetricService,
// 		di.services.hashService,
// 	)
// }

func (di *DI) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	di.workers.accrueNew.Start(ctx)
	time.Sleep(time.Second * 10)
	err := http.ListenAndServe(di.config.HTTP.Address, nil)

	di.infr.db.Close()

	return err
}
