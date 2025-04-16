package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	ssoGrpc "url-shortener/internal/clients/sso/grpc"
	"url-shortener/internal/config"
	remove "url-shortener/internal/http-server/handlers/delete"
	"url-shortener/internal/http-server/handlers/redirect"
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
	
	ssoClient, err := ssoGrpc.New(context.Background(),log,cfg.Clients.SSO.Address,cfg.Clients.SSO.Timeout,cfg.Clients.SSO.RetriesCount)
	if err != nil {
		log.Error("error init grpc client",sl.Err(err))
		os.Exit(1)
	}

	ssoClient.IsAdmin(context.Background(),1)
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

	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{cfg.HTTPServer.User: cfg.HTTPServer.Password}))

		r.Post("/save", save.New(log, storage))
		r.Delete("/{alias}", remove.New(log, storage))
	})

	router.Get("/{alias}", redirect.New(log, storage))

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

