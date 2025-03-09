package router

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"

	"github.com/gunzgo2mars/test-poke-service/app/internal/handler/auth"
	"github.com/gunzgo2mars/test-poke-service/app/internal/handler/info"
	"github.com/gunzgo2mars/test-poke-service/app/pkg/middleware"
)

type HttpServer struct {
	server *echo.Echo
	port   string

	// middleware
	jwt middleware.IJwt

	// Handlers
	authHandler auth.IAuthHandler
	infoHandler info.IInfoHandler
}

func New(
	port string,
	jwt middleware.IJwt,
	authHandler auth.IAuthHandler,
	infoHandler info.IInfoHandler,
) *HttpServer {
	e := echo.New()

	httpServer := &HttpServer{
		server:      e,
		port:        port,
		jwt:         jwt,
		authHandler: authHandler,
		infoHandler: infoHandler,
	}

	httpServer.initRouter()

	return httpServer
}

func (s *HttpServer) Start() error {
	return s.server.Start(fmt.Sprintf(":%s", s.port))
}

func (s *HttpServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *HttpServer) Server() *echo.Echo {
	return s.server
}
