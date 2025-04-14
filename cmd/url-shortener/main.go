package main

import (
	"fmt"
	"log/slog"
	"os"
	"url-shortener/internal/config"
	middlewareLogger "url-shortener/internal/http-server/middleware/logger"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage/sqllite"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// TODO: init config cleanenv for env

	cfg:=config.MustLoad()

	fmt.Println(cfg)

	// TODO: init logger slog
	log := config.SetupLogger(cfg.Env)	
	log.Info("Custom logger enabled in",slog.String("env",cfg.Env) )
	
	// TODO: init storage sqllite
	storage, err := sqllite.New(cfg.StoragePath)
	if err != nil {
		log.Error("error init storage",sl.Err(err))
		os.Exit(1)
	}

	_=storage

	// TODO: init router chi, render
	router := chi.NewRouter()

	// middleware
	router.Use(middleware.RequestID)
	router.Use(middlewareLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	// TODO: start server
}

