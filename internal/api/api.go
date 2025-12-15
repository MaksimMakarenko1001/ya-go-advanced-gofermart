package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/api/handler"
	createUserOrdersHandler "github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/api/handler/createUserOrders/v0"
	listUserOrdersHandler "github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/api/handler/listUserOrders/v0"
	createUserOrdersService "github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/service/createUserOrders/v0"
	listUserOrdersService "github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/internal/service/listUserOrders/v0"
)

type API struct {
	router *chi.Mux

	createUserOrdersService *createUserOrdersService.Service
	listUserOrdersService   *listUserOrdersService.Service
}

func New(
	// logger logger.HTTPLogger,
	createUserOrdersService *createUserOrdersService.Service,
	listUserOrdersService *listUserOrdersService.Service,
) *API {
	return &API{
		router: chi.NewRouter(),
		// logger:             logger,
		createUserOrdersService: createUserOrdersService,
		listUserOrdersService:   listUserOrdersService,
	}
}

func (api API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	api.router.ServeHTTP(w, r)
}

func (api API) HandleCreateUserOrders(middlewares ...handler.Middleware) {
	var h http.Handler = createUserOrdersHandler.Handle(api.createUserOrdersService.Do)

	middlewares = append(middlewares, handler.MiddlewareTypeContentTextPlain)
	h = handler.Conveyor(h, middlewares...)

	api.router.Post(createUserOrdersHandler.MethodPath, func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
	})
}

func (api API) HandleListUserOrders(middlewares ...handler.Middleware) {
	var h http.Handler = listUserOrdersHandler.Handle(api.listUserOrdersService.Do)

	h = handler.Conveyor(h, middlewares...)

	api.router.Get(listUserOrdersHandler.MethodPath, func(w http.ResponseWriter, r *http.Request) {
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
