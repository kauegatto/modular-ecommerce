package main

import (
	"ecommerce/Auth/Infra/config"
	"ecommerce/Auth/routes"
	"log"

	"github.com/gin-gonic/gin"
)

type App struct {
	router  *gin.Engine
	modules []Module
	config  *config.Config
}

type Module interface {
	RegisterRoutes(router *gin.RouterGroup)
	Name() string
}

func NewApp(config *config.Config) *App {
	router := gin.Default()
	return &App{
		router:  router,
		config:  config,
		modules: make([]Module, 0),
	}
}

func (a *App) RegisterModule(m Module) {
	moduleGroup := a.router.Group("/api/" + m.Name())

	m.RegisterRoutes(moduleGroup)
	a.modules = append(a.modules, m)
}

func (a *App) Start() error {
	return a.router.Run()
}

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	app := NewApp(config)

	authModule := routes.NewAuthHandler()
	app.RegisterModule(authModule)

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
