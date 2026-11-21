package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/thomasbuchinger/homelab-api/pkg/api"
	"github.com/thomasbuchinger/homelab-api/pkg/common"
	matrixhelper "github.com/thomasbuchinger/homelab-api/pkg/matrix-helper"
)

func main() {
	common.SetupViperConfig()
	serverConfig := common.GetServerConfig()
	logger := serverConfig.RootLogger.Named("MatrixHelper")
	router := api.SetupDefaultRouter()

	// Set up Signal handling
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	// Catch any Panics and crash the pod
	defer func() {
		if r := recover(); r != nil {
			logger.Errorf("matrix-helper failed: %v", r)
			logger.Sync()
			os.Exit(1)
		}
	}()
	defer stop()

	// Start HTTP Server
	go func() {
		logger.Infof("Starting matrix-helper on %s", serverConfig.BindAddr)
		if err := router.Run(serverConfig.BindAddr); err != nil {
			panic(fmt.Errorf("router: %w", err))
		}
	}()

	// Start SyncLoop
	matrixhelper.SyncLoop(ctx, logger)
	logger.Info("matrix-helper shut down") // We are only getting here if we exit without a panic
}
