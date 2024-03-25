package main

import (
	_ "github.com/lib/pq"
	"net/http"
	config "quests/configs"
	_ "quests/docs"
	"quests/handlers"
	slogpretty "quests/lib"
	"quests/repository"
	"quests/services"
)

// @title quests API
// @version 1.0
// @description API Server for Quests Application
// @host localhost:8081
// @securitydefinitions.basic BasicAuth
// @in header
// @name Authorization
func main() {
	//загружаем конфиг
	cfg := config.MustLoad()

	//инициализируем логгер
	logger := slogpretty.SetupLogger()
	logger.Info("Logger is start", "adress", cfg.HttpServer.Address)

	storage, err := repository.New(cfg.DbStorage)
	if err != nil {
		logger.Error("Database service is not start: ", err.Error())
		return
	}

	//создаем все необходимые таблицы
	err = storage.Init()
	if err != nil {
		logger.Error("Initialization database complete with error: ", err.Error())
	}
	logger.Info("Initialization database complete")

	repos := repository.NewRepository(storage.DB)
	services := services.NewService(repos)
	handlers := handlers.NewHandler(services)

	//запуск сервера
	server := &http.Server{
		Addr:         cfg.HttpServer.Address,
		Handler:      handlers.InitRoutes(),
		ReadTimeout:  cfg.HttpServer.TimeoutRequest,
		WriteTimeout: cfg.HttpServer.IdleTimeout,
		IdleTimeout:  cfg.HttpServer.IdleTimeout,
	}
	if err := server.ListenAndServe(); err != nil {
		logger.Error("Server does not started")
	}

}
