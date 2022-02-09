package api

import (
	"net/http"

	"channels-instagram-dm/domain"
	"github.com/gorilla/mux"
)

func InitRoutes(ctx domain.RuntimeContext, r *mux.Router) {
	RouteHandler(ctx, r, "/health", HealthCheck).Methods(http.MethodGet)
	RouteHandler(ctx, r, "/slots", Slots).Methods(http.MethodGet)
	RouteHandler(ctx, r, "/slots/refresh", RefreshSlots).Methods(http.MethodPost)

	RouteHandler(ctx, r, "/account/all", GetAllAccounts).Methods(http.MethodGet)
	RouteHandler(ctx, r, "/account/{external_id}", GetAccount).Methods(http.MethodGet)
	RouteHandler(ctx, r, "/account", AddAccount).Methods(http.MethodPost)
	RouteHandler(ctx, r, "/account/{external_id}", DeleteAccount).Methods(http.MethodDelete)

	RouteHandler(ctx, r, "/login", Login).Methods(http.MethodPost)
	RouteHandler(ctx, r, "/logout/{external_id}", Logout).Methods(http.MethodPost)

	RouteHandler(ctx, r, "/account/resume/{external_id}", ResumeAccount).Methods(http.MethodPost)
	RouteHandler(ctx, r, "/account/suspend/{external_id}", SuspendAccount).Methods(http.MethodPost)

	RouteHandler(ctx, r, "/account/activity/{external_id}", GetActivityLog).Methods(http.MethodGet)
}
