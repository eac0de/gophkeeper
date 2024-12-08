package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"auth/internal/api/handlers"
	"auth/internal/api/inmiddlewares"
	"auth/internal/config"
	"auth/internal/grpcserver"
	"auth/internal/services"
	"auth/internal/storage/psql"
	"auth/pkg/emailsender"

	"github.com/gin-gonic/gin"
)

func setupRouter(
	sessionService *services.SessionService,
	authService *services.AuthService,
) *gin.Engine {
	router := gin.Default()
	rootGroup := router.Group("api/auth")

	authHandlers := handlers.NewAuthHandlers(
		authService,
		sessionService,
	)

	rootGroup.POST("code/generate/", authHandlers.GenerateEmailCodeHandler)
	rootGroup.POST("code/check/", authHandlers.NewCheckEmailCodeHandler("/api/auth/token/"))
	rootGroup.POST("token/", authHandlers.NewRefreshTokenHandler("/api/auth/token/"))
	rootGroup.DELETE("token/", authHandlers.NewDeleteCurrentSession("/api/auth/token/"))

	authenticatedGroup := rootGroup.Group("/", inmiddlewares.NewAuthMiddleware(sessionService))
	authenticatedGroup.GET("/sessions/", authHandlers.GetUserSessionsHandler)
	authenticatedGroup.DELETE("/sessions/:id/", authHandlers.DeleteSession)
	return router
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.MustLoad()

	if !cfg.IsDev {
		gin.SetMode(gin.ReleaseMode)
	}

	psqlStorage, err := psql.New(
		ctx,
		cfg.PSQLHost,
		cfg.PSQLPort,
		cfg.PSQLUsername,
		cfg.PSQLPassword,
		cfg.PSQLDBName,
	)
	if err != nil {
		panic(err)
	}
	defer psqlStorage.Close()

	var emailSender emailsender.IEmailSender
	if cfg.IsDev {
		emailSender = emailsender.NewMock()
	} else {
		emailSender = emailsender.New(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUsername, cfg.SMTPPassword)
	}

	sessionService := services.NewSessionService(cfg.JWTSecretKey, cfg.JWTAccessExp, cfg.JWTRefreshExp, psqlStorage)
	authService := services.NewAuthService(psqlStorage, emailSender)

	gprcAuthServer := grpcserver.NewAuthGRPCServer(cfg.GPRCServerAddress, sessionService)
	go gprcAuthServer.Run()

	r := setupRouter(sessionService, authService)
	go r.Run(cfg.ServerAddress)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-sigChan
}
