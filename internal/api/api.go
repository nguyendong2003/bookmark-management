package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nguyendong2003/bookmark-management/internal/handler"
	"github.com/nguyendong2003/bookmark-management/internal/service"
)

type Engine interface {
	Start() error
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type api struct {
	app *gin.Engine
	cfg *Config
}

func New(cfg *Config) Engine {
	a := &api{
		app: gin.New(),
		cfg: cfg,
	}

	a.registerEndpoint()

	return a
}

func (a *api) Start() error {
	return a.app.Run(fmt.Sprintf(":%s", a.cfg.AppPort))
}

func (a *api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.app.ServeHTTP(w, r)
}

func (a *api) registerEndpoint() {
	passwordService := service.NewPassword()
	passwordHandler := handler.NewPassword(passwordService)
	a.app.GET("/gen-pass", passwordHandler.GenPass)
}
