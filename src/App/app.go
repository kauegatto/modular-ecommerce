package app

import (
	config "ecommerce/SharedKernel"

	"github.com/gin-gonic/gin"
)

type App struct {
	router  *gin.Engine
	modules []Module
	config  *config.Configuration
}

type Module interface {
	RegisterRoutes(router *gin.RouterGroup)
	Name() string
}

func NewApp(config *config.Configuration) *App {
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
