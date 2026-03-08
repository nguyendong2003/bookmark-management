package api

import (
	"net/http"

	"github.com/nguyendong2003/bookmark-management/internal/handler"
	"github.com/nguyendong2003/bookmark-management/internal/service"
)

type apiMux struct {
	mux *http.ServeMux
}

func NewMux() Engine {
	a := &apiMux{
		mux: http.NewServeMux(),
	}

	a.registerEndpoint()

	return a
}

func (a *apiMux) Start() error {
	return http.ListenAndServe(":8080", a.mux)
}

func (a *apiMux) registerEndpoint() {
	passwordService := service.NewPassword()
	passwordHandler := handler.NewPassword(passwordService)
	a.mux.HandleFunc("/gen-pass-mux", passwordHandler.GenPassForMux)
}
