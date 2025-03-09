package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	libMiddleware "github.com/gunzgo2mars/test-poke-service/app/pkg/middleware"
)

func (s *HttpServer) initRouter() {
	root := s.server

	root.Use(middleware.CORS())
	root.Use(middleware.RateLimiterWithConfig(libMiddleware.RateLimiterConfig()))

	root.GET("/checkhealth", func(c echo.Context) error {
		return c.JSON(
			http.StatusOK,
			"ok",
		)
	})

	// Public endpoints.
	public := root.Group("/api/v1")
	public.POST("/register", s.authHandler.Register)
	public.POST("/login", s.authHandler.Login)

	// Private endpoints.
	private := root.Group("/api/v1", s.jwt.ValidateMiddleware)
	private.GET("/pokemon/:name", s.infoHandler.GetPokemonByName)
	private.GET("/pokemon/:name/ability", s.infoHandler.GetPokemonAbilities)
	private.GET("/pokemon/random", s.infoHandler.Random)
}
