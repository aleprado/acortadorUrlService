package main

import (
	"acortadorUrlService/components/config"
	"acortadorUrlService/components/database"
	"acortadorUrlService/components/logger"
	"acortadorUrlService/components/metrics"
	"acortadorUrlService/url-api/controller"
	"acortadorUrlService/url-api/service"
	"acortadorUrlService/url-api/web"
	"context"
	"net/http"
	"os"
)

func main() {
	logger.Init()
	ctx := context.Background()
	cfg := config.LoadConfig()
	metrics.Init(cfg.Region)
	dbClient, err := database.NewDDBClient(ctx, cfg)
	if err != nil {
		logger.LogError("No se pudo inicializar la base de datos", "error", err)
		os.Exit(1)
	}
	
	urlShortener := service.NewUrlShortener(dbClient)
	urlController := controller.NewUrlController(urlShortener, cfg)

	router := web.NewHttpHandler("v1")
	urlController.MountIn(router)
	port := cfg.Port
	logger.LogInfo("Starting server on port: " + port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		logger.LogError("Server failed", "error", err)
	}
}
