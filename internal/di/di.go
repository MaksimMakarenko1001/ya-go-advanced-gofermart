package di

import (
	"context"
	"log"
	"net/http"

	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/api"
	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/api/handler"
	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/db"
	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/repository/accrual"
	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/repository/hash"
	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/repository/jwt"
	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/repository/lock"
	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/repository/order"
	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/repository/user"
	accrueNewOrders "github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/service/accrueNewOrders/v0"
	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/service/auth"
	getAccrualInfoByOrders "github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/service/getAccrualInfoByOrders/v0"
	getUserBalance "github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/service/getUserBalance/v0"
	getUserOrders "github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/service/getUserOrders/v0"
	getUserWithdrawals "github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/service/getUserWithdrawals/v0"
	postUserBalanceWithdraw "github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/service/postUserBalanceWithdraw/v0"
	postUserLogin "github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/service/postUserLogin/v0"
	postUserOrders "github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/service/postUserOrders/v0"
	postUserRegister "github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/service/postUserRegister/v0"
	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/worker"
	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/pkg/backoff"
)

type DI struct {
	config       *diConfig
	repositories struct {
		accrual *accrual.Repository
		order   *order.Repository
		user    *user.Repository
		lock    *lock.Repository
		hash    *hash.Repository
		jwt     *jwt.Repository
	}
	services struct {
		included struct {
			getAccrualInfoByOrdersService *getAccrualInfoByOrders.Service
		}
		authService                    *auth.Service
		accrueNewOrdersService         *accrueNewOrders.Service
		postUserOrdersService          *postUserOrders.Service
		getUserOrdersService           *getUserOrders.Service
		getUserBalanceService          *getUserBalance.Service
		postUserBalanceWithdrawService *postUserBalanceWithdraw.Service
		getUserWithdrawalsService      *getUserWithdrawals.Service
		postUserRegisterService        *postUserRegister.Service
		postUserLoginService           *postUserLogin.Service
	}
	workers struct {
		accrueNew        *worker.Worker
		accrueProcessing *worker.Worker
	}
	api struct {
		external *api.API
	}
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
	di.initAPI()
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
		backoff.NewBackoff(accrual.ClassifyHTTPError, di.config.Repository.Accrual.MaxRetries, di.config.Repository.Accrual.MinDelay),
		di.config.Repository.Accrual.DeltaDelay,
	))
	di.repositories.lock = lock.New(di.infr.db)
	di.repositories.order = order.New(di.infr.db)
	di.repositories.user = user.New(di.infr.db)
	di.repositories.hash = hash.New(di.config.Repository.Hash)
	di.repositories.jwt = jwt.New(di.config.Repository.Jwt)
}

func (di *DI) initServices() {
	di.services.included.getAccrualInfoByOrdersService = getAccrualInfoByOrders.New(di.config.Service.GetAccrualInfoByOrders, di.repositories.accrual)
	di.services.authService = auth.New(di.repositories.jwt)

	di.services.accrueNewOrdersService = accrueNewOrders.New(di.config.Service.AccrueNewOrders, di.repositories.order, di.services.included.getAccrualInfoByOrdersService)
	di.services.postUserOrdersService = postUserOrders.New(di.repositories.order, di.repositories.jwt)
	di.services.getUserOrdersService = getUserOrders.New(di.repositories.order, di.repositories.jwt)
	di.services.getUserBalanceService = getUserBalance.New(di.repositories.order, di.repositories.jwt)
	di.services.postUserBalanceWithdrawService = postUserBalanceWithdraw.New(di.repositories.order, di.repositories.jwt)
	di.services.getUserWithdrawalsService = getUserWithdrawals.New(di.repositories.order, di.repositories.jwt)
	di.services.postUserRegisterService = postUserRegister.New(di.repositories.user, di.repositories.hash, di.repositories.jwt)
	di.services.postUserLoginService = postUserLogin.New(di.repositories.user, di.repositories.hash, di.repositories.jwt)

}

func (di *DI) initWorkers() {
	di.workers.accrueNew = worker.New(
		di.config.Worker.AccrueNew,
		di.repositories.lock,
		di.config.AppName,
		"accrue_new",
		"",
		di.services.accrueNewOrdersService.Do,
	)
}

func (di *DI) initAPI() {
	di.api.external = api.New(
		// logger.New(di.config.Logger),
		di.services.authService,
		di.services.postUserOrdersService,
		di.services.getUserOrdersService,
		di.services.getUserBalanceService,
		di.services.postUserBalanceWithdrawService,
		di.services.getUserWithdrawalsService,
		di.services.postUserRegisterService,
		di.services.postUserLoginService,
	)

	di.api.external.HandlePostUserOrders(di.api.external.WithJwtAuth, handler.MiddlewareTypeContentTextPlain)
	di.api.external.HandleGetUserOrders(di.api.external.WithJwtAuth)
	di.api.external.HandleGetUserBalance(di.api.external.WithJwtAuth)
	di.api.external.HandlePostUserBalanceWithdraw(di.api.external.WithJwtAuth, handler.MiddlewareTypeContentApplicationJSON)
	di.api.external.HandleGetUserWithdrawals(di.api.external.WithJwtAuth)

	di.api.external.HandlePostUserRegister(handler.MiddlewareTypeContentApplicationJSON)
	di.api.external.HandlePostUserLogin(handler.MiddlewareTypeContentApplicationJSON)
}

func (di *DI) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	di.workers.accrueNew.Start(ctx)

	err := http.ListenAndServe(di.config.HTTP.Address, handler.Conveyor(
		di.api.external,
		handler.MiddlewareCompress,
	))

	di.infr.db.Close()

	return err
}
