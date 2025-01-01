package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/eac0de/gophkeeper/auth/pkg/outmiddlewares"
	"github.com/eac0de/gophkeeper/internal/api/handlers"
	"github.com/eac0de/gophkeeper/internal/config"
	"github.com/eac0de/gophkeeper/internal/services"
	"github.com/eac0de/gophkeeper/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/gin-gonic/gin"
)

func setupRouter(
	authServiceConn *grpc.ClientConn,
	userDataService *services.UserDataService,
) *gin.Engine {
	router := gin.Default()
	rootGroup := router.Group("api/gophkeeper/")
	authenticatedGroup := rootGroup.Group("/", outmiddlewares.NewAuthMiddleware(authServiceConn))

	userDataHandlers := handlers.NewUserDataHandlers(userDataService)

	authenticatedGroup.GET("/auth_info/", userDataHandlers.GetUserAuthInfoList)
	authenticatedGroup.GET("/auth_info/:id/", userDataHandlers.GetUserAuthInfo)
	authenticatedGroup.DELETE("/auth_info/:id/", userDataHandlers.DeleteUserAuthInfo)
	authenticatedGroup.PUT("/auth_info/:id/", userDataHandlers.UpdateUserAuthInfo)
	authenticatedGroup.POST("/auth_info/", userDataHandlers.InsertUserAuthInfo)

	authenticatedGroup.GET("/text_data/", userDataHandlers.GetUserTextDataList)
	authenticatedGroup.GET("/text_data/:id/", userDataHandlers.GetUserTextData)
	authenticatedGroup.DELETE("/text_data/:id/", userDataHandlers.DeleteUserTextData)
	authenticatedGroup.PUT("/text_data/:id/", userDataHandlers.UpdateUserTextData)
	authenticatedGroup.POST("/text_data/", userDataHandlers.InsertUserTextData)

	authenticatedGroup.GET("/file_data/", userDataHandlers.GetUserFileDataList)
	authenticatedGroup.GET("/file_data/:id/", userDataHandlers.GetUserFileData)
	authenticatedGroup.GET("/file_data/:id/download/", userDataHandlers.DownloadUserFile)
	authenticatedGroup.DELETE("/file_data/:id/", userDataHandlers.DeleteUserFileData)
	authenticatedGroup.PUT("/file_data/:id/", userDataHandlers.UpdateUserFileData)
	authenticatedGroup.POST("/file_data/", userDataHandlers.InsertUserFileData)

	authenticatedGroup.GET("bank_cards/", userDataHandlers.GetUserBankCardList)
	authenticatedGroup.GET("bank_cards/:id/", userDataHandlers.GetUserBankCard)
	authenticatedGroup.DELETE("bank_cards/:id/", userDataHandlers.DeleteUserBankCard)
	authenticatedGroup.PUT("bank_cards/:id/", userDataHandlers.UpdateUserBankCard)
	authenticatedGroup.POST("bank_cards/", userDataHandlers.InsertUserBankCard)

	return router
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.MustLoad()

	if !cfg.IsDev {
		gin.SetMode(gin.ReleaseMode)
	}

	gophKeeperStorage, err := storage.NewGophKeeperStorage(
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
	err = gophKeeperStorage.Migrate(ctx, "./migrations", false)
	if err != nil {
		panic(err)
	}
	defer gophKeeperStorage.Close()

	userDataService := services.NewUserDataService(gophKeeperStorage)
	authServiceConn, err := grpc.NewClient(cfg.AuthGRPCServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	r := setupRouter(authServiceConn, userDataService)
	go r.Run(cfg.ServerAddress)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-sigChan
}
