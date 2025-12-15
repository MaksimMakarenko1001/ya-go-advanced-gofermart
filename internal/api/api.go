package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/api/handler"
	getUserBalanceHandler "github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/api/handler/getUserBalance/v0"
	getUserOrdersHandler "github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/api/handler/getUserOrders/v0"
	postUserBalanceWithdrawHandler "github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/api/handler/postUserBalanceWithdraw/v0"
	postUserOrdersHandler "github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/api/handler/postUserOrders/v0"
	getUserBalanceService "github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/service/getUserBalance/v0"
	getUserOrdersService "github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/service/getUserOrders/v0"
	postUserBalanceWithdrawService "github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/service/postUserBalanceWithdraw/v0"
	postUserOrdersService "github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/service/postUserOrders/v0"
)

type API struct {
	router *chi.Mux

	postUserOrdersService          *postUserOrdersService.Service
	getUserOrdersService           *getUserOrdersService.Service
	getUserBalanceService          *getUserBalanceService.Service
	postUserBalanceWithdrawService *postUserBalanceWithdrawService.Service
}

func New(
	// logger logger.HTTPLogger,
	postUserOrdersService *postUserOrdersService.Service,
	getUserOrdersService *getUserOrdersService.Service,
	getUserBalanceService *getUserBalanceService.Service,
	postUserBalanceWithdrawService *postUserBalanceWithdrawService.Service,
) *API {
	return &API{
		router: chi.NewRouter(),
		// logger:             logger,
		postUserOrdersService:          postUserOrdersService,
		getUserOrdersService:           getUserOrdersService,
		getUserBalanceService:          getUserBalanceService,
		postUserBalanceWithdrawService: postUserBalanceWithdrawService,
	}
}

func (api API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	api.router.ServeHTTP(w, r)
}

func (api API) HandlePostUserOrders(middlewares ...handler.Middleware) {
	var h http.Handler = postUserOrdersHandler.Handle(api.postUserOrdersService.Do)

	middlewares = append(middlewares, handler.MiddlewareTypeContentTextPlain)
	h = handler.Conveyor(h, middlewares...)

	api.router.Post(postUserOrdersHandler.MethodPath, func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
}

func (api API) HandleListUserOrders(middlewares ...handler.Middleware) {
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

	middlewares = append(middlewares, handler.MiddlewareTypeContentApplicationJSON)
	h = handler.Conveyor(h, middlewares...)

	api.router.Post(postUserBalanceWithdrawHandler.MethodPath, func(w http.ResponseWriter, r *http.Request) {
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
