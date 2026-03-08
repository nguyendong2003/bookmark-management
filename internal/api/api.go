package api

import (
	"github.com/gin-gonic/gin"
	"github.com/nguyendong2003/bookmark-management/internal/handler"
	"github.com/nguyendong2003/bookmark-management/internal/service"
)

type Engine interface {
	Start() error
}

type api struct {
	app *gin.Engine
}

func New() Engine {
	a := &api{
		app: gin.New(),
	}

	a.registerEndpoint()

	return a
}

func (a *api) Start() error {
	return a.app.Run(":8080")
}

func (a *api) registerEndpoint() {
	passwordService := service.NewPassword()
	passwordHandler := handler.NewPassword(passwordService)
	a.app.GET("/gen-pass", passwordHandler.GenPass)
}
