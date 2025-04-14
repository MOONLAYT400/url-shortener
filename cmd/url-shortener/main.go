package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/http-server/handlers/url/save"
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

	router.Post("/save", save.New(log, storage))

	// TODO: start server
	log.Info("server started",slog.String("addr",cfg.HTTPServer.Address))
	srv := http.Server{
		Addr:    cfg.HTTPServer.Address,
		Handler: router,
		ReadTimeout: cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout: cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("error start server",sl.Err(err))
		os.Exit(1)
	}

	log.Error("server emergency stop")
}

