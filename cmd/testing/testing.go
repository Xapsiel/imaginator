package main

import (
	"flag"
	"log/slog"
	"os"

	"imageBot/internal/api"
	"imageBot/internal/config"
)

func main() {
	configPath := flag.String("c", "configs/config.yaml", "The path to the configuration file")
	flag.Parse()
	cfg, err := config.New(*configPath)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	T2I_api := api.Text2ImageAPI{
		URL:    cfg.FB.URL,
		APIKey: cfg.FB.APIKey,
		Secret: cfg.FB.Secret,
	}
	bytea, err := T2I_api.Draw("Даня гофер", 1024, 1024)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	os.WriteFile("pic.jpg", bytea, 0644)

}
