package main

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/nikitapozdeev/feed-repost-bot/internal/app"
	"github.com/nikitapozdeev/feed-repost-bot/internal/clients/vk"
	"github.com/nikitapozdeev/feed-repost-bot/internal/config"
	"github.com/nikitapozdeev/feed-repost-bot/internal/db"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stdout, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	var cfg config.Config
	err := cleanenv.ReadConfig("./config.yml", &cfg)
	if err != nil {
		return fmt.Errorf("Failed to read configuration: %w", err)
	}

	db, err := db.NewDB("store")
	if err != nil {
		return fmt.Errorf("Failed to create database: %w", err)
	}
	defer db.Close()

	// create clients vk, youtube, facebook instagram, etc.
	vkClient := vk.NewClient(
		cfg.VK.Host,
		cfg.VK.BasePath,
		cfg.VK.Version,
		cfg.VK.Token,
	)

	a, err := app.NewApp(&cfg, db, vkClient)
	if err != nil {
		return fmt.Errorf("Failed to create app: %w", err)
	}
	a.Run()

	return nil
}
