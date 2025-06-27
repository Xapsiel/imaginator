package main

import (
	"flag"
	"log/slog"

	"imageBot/internal/api"
	"imageBot/internal/bot"
	"imageBot/internal/config"
	"imageBot/internal/repository"
	"imageBot/internal/service"
)

func main() {
	configPath := flag.String("c", "configs/config.yaml", "The path to the configuration file")
	flag.Parse()
	cfg, err := config.New(*configPath)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	db, err := repository.NewPostgresDB(cfg.DatabaseConfig)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	slog.Info("Initializing db")
	repo := repository.NewRepository(db)
	slog.Info("Initializing repos")

	api_fb := api.New(cfg.FB)
	services := service.NewService(repo, api_fb)
	bot := bot.New(cfg.Bot, cfg.Prompt, services)
	bot.Start()
}
