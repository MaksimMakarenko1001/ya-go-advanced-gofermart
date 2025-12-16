package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/api/handler"
	getUserBalanceHandler "github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/api/handler/getUserBalance/v0"
	getUserOrdersHandler "github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/api/handler/getUserOrders/v0"
	getUserWithdrawalsHandler "github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/api/handler/getUserWithdrawals/v0"
	postUserBalanceWithdrawHandler "github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/api/handler/postUserBalanceWithdraw/v0"
	postUserLoginHandler "github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/api/handler/postUserLogin/v0"
	postUserOrdersHandler "github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/api/handler/postUserOrders/v0"
	postUserRegisterHandler "github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/api/handler/postUserRegister/v0"
	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/service/auth"
	getUserBalanceService "github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/service/getUserBalance/v0"
	getUserOrdersService "github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/service/getUserOrders/v0"
	getUserWithdrawalsService "github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/service/getUserWithdrawals/v0"
	postUserBalanceWithdrawService "github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/service/postUserBalanceWithdraw/v0"
	postUserLoginService "github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/service/postUserLogin/v0"
	postUserOrdersService "github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/service/postUserOrders/v0"
	postUserRegisterService "github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/service/postUserRegister/v0"
	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/pkg"
)

type API struct {
	router *chi.Mux

	authService                    *auth.Service
	postUserOrdersService          *postUserOrdersService.Service
	getUserOrdersService           *getUserOrdersService.Service
	getUserBalanceService          *getUserBalanceService.Service
	postUserBalanceWithdrawService *postUserBalanceWithdrawService.Service
	getUserWithdrawalsService      *getUserWithdrawalsService.Service
	postUserRegisterService        *postUserRegisterService.Service
	postUserLoginService           *postUserLoginService.Service
}

func New(
	// logger logger.HTTPLogger,
	authService *auth.Service,
	postUserOrdersService *postUserOrdersService.Service,
	getUserOrdersService *getUserOrdersService.Service,
	getUserBalanceService *getUserBalanceService.Service,
	postUserBalanceWithdrawService *postUserBalanceWithdrawService.Service,
	getUserWithdrawalsService *getUserWithdrawalsService.Service,
	postUserRegisterService *postUserRegisterService.Service,
	postUserLoginService *postUserLoginService.Service,
) *API {
	return &API{
		router: chi.NewRouter(),
		// logger:             logger,
		authService:                    authService,
		postUserOrdersService:          postUserOrdersService,
		getUserOrdersService:           getUserOrdersService,
		getUserBalanceService:          getUserBalanceService,
		postUserBalanceWithdrawService: postUserBalanceWithdrawService,
		getUserWithdrawalsService:      getUserWithdrawalsService,
		postUserRegisterService:        postUserRegisterService,
		postUserLoginService:           postUserLoginService,
	}
}

func (api API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	api.router.ServeHTTP(w, r)
}

func (api API) HandlePostUserOrders(middlewares ...handler.Middleware) {
	var h http.Handler = postUserOrdersHandler.Handle(api.postUserOrdersService.Do)

	h = handler.Conveyor(h, middlewares...)

	api.router.Post(postUserOrdersHandler.MethodPath, func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
}

func (api API) HandleGetUserOrders(middlewares ...handler.Middleware) {
	var h http.Handler = getUserOrdersHandler.Handle(api.getUserOrdersService.Do)

	h = handler.Conveyor(h, middlewares...)

	api.router.Get(getUserOrdersHandler.MethodPath, func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
}

func (api API) HandleGetUserBalance(middlewares ...handler.Middleware) {
	var h http.Handler = getUserBalanceHandler.Handle(api.getUserBalanceService.Do)

	h = handler.Conveyor(h, middlewares...)

	api.router.Get(getUserBalanceHandler.MethodPath, func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
}

func (api API) HandlePostUserBalanceWithdraw(middlewares ...handler.Middleware) {
	var h http.Handler = postUserBalanceWithdrawHandler.Handle(api.postUserBalanceWithdrawService.Do)

	h = handler.Conveyor(h, middlewares...)

	api.router.Post(postUserBalanceWithdrawHandler.MethodPath, func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
}

func (api API) HandleGetUserWithdrawals(middlewares ...handler.Middleware) {
	var h http.Handler = getUserWithdrawalsHandler.Handle(api.getUserWithdrawalsService.Do)

	h = handler.Conveyor(h, middlewares...)

	api.router.Get(getUserWithdrawalsHandler.MethodPath, func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
}

func (api API) HandlePostUserRegister(middlewares ...handler.Middleware) {
	var h http.Handler = postUserRegisterHandler.Handle(api.postUserRegisterService.Do)

	h = handler.Conveyor(h, middlewares...)

	api.router.Post(postUserRegisterHandler.MethodPath, func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
}

func (api API) HandlePostUserLogin(middlewares ...handler.Middleware) {
	var h http.Handler = postUserLoginHandler.Handle(api.postUserLoginService.Do)

	h = handler.Conveyor(h, middlewares...)

	api.router.Post(postUserLoginHandler.MethodPath, func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
}

func (api API) WithJwtAuth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken := ""
		if authHeader := r.Header.Get("Authorization"); authHeader != "" {
			if !strings.HasPrefix(authHeader, "Bearer ") {
				handler.WriteError(w, fmt.Errorf("invalid authorization header format, %w", pkg.ErrBadRequest))
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")
			ok, err := api.authService.ValidateToken(r.Context(), token)
			if err != nil {
				handler.WriteError(w, err)
				return
			}
			if !ok {
				handler.WriteError(w, pkg.ErrUnauthorized)
				return
			}
			accessToken = token
		}

		r.Header.Set(handler.HeaderAccessToken, accessToken)
		h.ServeHTTP(w, r)
	})
}

// func (api API) WithLogging(h http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		start := time.Now()

// 		var resp ResponseInfo
// 		rw := responseWriter{
// 			ResponseWriter: w,
// 			response:       &resp,
// 		}

// 		h.ServeHTTP(&rw, r)

// 		api.logger.LogHTTP(logger.HTTPInfo{
// 			URI:      r.RequestURI,
// 			Method:   r.Method,
// 			Duration: time.Since(start),
// 			Response: logger.ResponseInfo{
// 				Size:   resp.Size,
// 				Status: resp.Status,
// 				Body:   resp.Body,
// 			},
// 		})
// 	})
// }
