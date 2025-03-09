package main

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	authSvc "github.com/gunzgo2mars/test-poke-service/app/internal/core/service/auth"
	pokemonSvc "github.com/gunzgo2mars/test-poke-service/app/internal/core/service/pokemon"
	authHandler "github.com/gunzgo2mars/test-poke-service/app/internal/handler/auth"
	infoHandler "github.com/gunzgo2mars/test-poke-service/app/internal/handler/info"
	extRepo "github.com/gunzgo2mars/test-poke-service/app/internal/repository/ext"
	usersRepo "github.com/gunzgo2mars/test-poke-service/app/internal/repository/users"
	"github.com/gunzgo2mars/test-poke-service/app/internal/router"
	"github.com/gunzgo2mars/test-poke-service/app/pkg/cache"
	"github.com/gunzgo2mars/test-poke-service/app/pkg/configurer"
	"github.com/gunzgo2mars/test-poke-service/app/pkg/database"
	"github.com/gunzgo2mars/test-poke-service/app/pkg/logger"
	"github.com/gunzgo2mars/test-poke-service/app/pkg/middleware"
	"github.com/gunzgo2mars/test-poke-service/app/pkg/monitor"
	"github.com/gunzgo2mars/test-poke-service/app/pkg/resty"
	"github.com/gunzgo2mars/test-poke-service/app/pkg/utils"
)

func main() {
	// initializing app context and config
	ctx, cancel := context.WithCancel(context.Background())
	var conf *configurer.AppConfig

	if err := configurer.LoadConfig(&conf, "config", "config", "yaml", "APPENV"); err != nil {
		panic(errors.Wrap(err, "Load Config"))
	}

	if err := configurer.LoadDotEnv(&conf.Secrets, "./config/", "SECRET", "APPENV"); err != nil {
		panic(errors.Wrap(err, "Load DotEnv"))
	}
	fmt.Printf("Debug Conf: %v \n", conf)

	// initializing dependencies.
	logger.InitLogger(conf.Log.Env)
	defer logger.Sync()

	jwtMiddleware := middleware.NewJwt(conf.Secrets.JwtKey)
	validator := utils.NewValidator()
	resty := resty.New()

	mysql := database.NewMysql(conf)
	db, err := mysql.Connect()
	if err != nil {
		logger.Error(
			ctx,
			fmt.Sprintf("Error: %s", err.Error()),
			zap.String("type", "exec"),
		)
	}
	defer database.Close(db)
	logger.Info(
		ctx,
		"MySQL Database connected!",
		zap.String("type", "exec"),
	)

	cache := cache.NewRedis(
		&redis.Options{Addr: fmt.Sprintf("%s:%s", conf.Redis.Address, conf.Redis.Port)},
	)

	// initializing repositories.
	userRepo := usersRepo.New(db)
	extRepo := extRepo.New(resty, cache, conf)

	// initializing services.
	authSvc := authSvc.New(jwtMiddleware, userRepo)
	pokemonSvc := pokemonSvc.New(extRepo)

	// initializing handlers.
	authHandler := authHandler.New(authSvc, validator)
	infoHandler := infoHandler.New(pokemonSvc)

	// initializing Http Server.
	httpServer := router.New(
		conf.App.Port,
		jwtMiddleware,
		authHandler,
		infoHandler,
	)
	go func() {
		if err := httpServer.Start(); err != nil {
			logger.Fatal(
				ctx,
				fmt.Sprintf("Error: %s", err.Error()),
				zap.String("type", "exec"),
			)
		}
	}()

	// Graceful shutdown http server.
	monitor.GracefulShutdownHttpServer(ctx, cancel, httpServer.Server())
}
