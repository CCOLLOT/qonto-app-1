// boilerplate app code with graceful shutdown, an HTTP server struct with basic builder functions (routes, middlewares etc.)
package app

import (
	"context"
	"net/http"
	"os"

	"github.com/CCOLLOT/appnametochange/internal/server"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	MESSAGE = "This application is running smoothly"
)

type App struct {
	log  *zap.Logger
	http *server.HTTPServer
}

func New(log *zap.Logger) (*App, error) {
	srv, err := server.New(log)
	if err != nil {
		return nil, err
	}
	return &App{
		log:  log,
		http: srv,
	}, nil
}

func (app *App) Name() string {
	return "appnametochange"
}

func (app *App) Start() error {
	app.http.Engine.GET("/healthz", healthCheck)
	go func() {
		err := app.http.Server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			app.log.Error(err.Error())
			os.Exit(2)
		}
	}()
	return nil
}

func healthCheck(ctx *gin.Context) {
	ctx.JSON(200, MESSAGE)
}

func (app *App) Shutdown() error {
	app.log.Info("shutting down HTTP server...")
	if err := app.http.Server.Shutdown(context.Background()); err != nil {
		app.log.Error(err.Error())
		return err
	}
	return nil
}