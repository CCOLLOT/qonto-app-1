package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type HTTPServer struct {
	Server *http.Server
	Engine *gin.Engine
	Log    *zap.Logger
}

func New(log *zap.Logger) (*HTTPServer, error) {
	eng := gin.New()
	addr := fmt.Sprintf("%s:%d", "0.0.0.0", 8080)
	return &HTTPServer{
		Log:    log,
		Engine: eng,
		Server: &http.Server{
			Addr: addr,
			Handler: eng,
		},
	}, nil
}
