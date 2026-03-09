package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/nguyendong2003/bookmark-management/docs"
	"github.com/nguyendong2003/bookmark-management/internal/handler"
	"github.com/nguyendong2003/bookmark-management/internal/repository"
	"github.com/nguyendong2003/bookmark-management/internal/service"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Engine interface {
	Start() error
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type engine struct {
	app         *gin.Engine
	cfg         *Config
	redisClient *redis.Client
}

func NewEngine(cfg *Config, redisClient *redis.Client) Engine {
	e := &engine{
		app:         gin.New(),
		cfg:         cfg,
		redisClient: redisClient,
	}

	e.initRoutes()

	return e
}

func (e *engine) Start() error {
	return e.app.Run(fmt.Sprintf(":%s", e.cfg.AppPort))
}

func (e *engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e.app.ServeHTTP(w, r)
}

func (e *engine) initRoutes() {
	// create handlers
	passwordService := service.NewPassword()
	passwordHandler := handler.NewPassword(passwordService)

	// urlstorage handler
	urlStorage := repository.NewURLStorage(e.redisClient)
	shortenURLService := service.NewShortenURL(urlStorage, passwordService)
	shortenURLHandler := handler.NewShortenURL(shortenURLService)

	// register handlers to endpoints
	e.app.GET("/gen-pass", passwordHandler.GenPass)
	e.app.POST("/shorten", shortenURLHandler.ShortenURL)
	e.app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
