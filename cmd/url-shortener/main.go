package main

import (
	"fmt"
	"log/slog"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage/sqllite"
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

	// TODO: start server
}

