package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(port string, router *gin.Engine) *Server {
	addr := fmt.Sprintf(":%s", port)

	return &Server{
		httpServer: &http.Server{
			Addr:    addr,
			Handler: router,
		},
	}
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) GetHTTPServer() *http.Server {
	return s.httpServer
}
