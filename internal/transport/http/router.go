package transport

import (
	"github.com/gorilla/mux"

	"weaver-edge/internal/transport/http/handler"
	"weaver-edge/internal/transport/http/middleware"
)

func NewRouter(handler *handler.Handler, contextSecret string) *mux.Router {
	r := mux.NewRouter()

	r.Use(middleware.MiddlewareLogging)
	r.HandleFunc("/health", handler.Health).Methods("GET")

	r.HandleFunc("/webhooks/telegram", handler.TelegramWebhook).Methods("POST")

	contextRouter := r.PathPrefix("/context").Subrouter()
	contextRouter.HandleFunc("", handler.Context).Methods("POST")
	contextRouter.Use(middleware.AuthMiddleware(contextSecret))

	return r
}
