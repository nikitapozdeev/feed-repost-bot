package main

import (
	"fmt"
	"os"
	"syscall"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/nikitapozdeev/feed-repost-bot/internal/app"
	"github.com/nikitapozdeev/feed-repost-bot/internal/clients/vk"
	"github.com/nikitapozdeev/feed-repost-bot/internal/config"
	"github.com/nikitapozdeev/feed-repost-bot/internal/storage/sqldb"
	"github.com/nikitapozdeev/feed-repost-bot/pkg/shutdown"
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

	storage, err := sqldb.NewDB("store")
	if err != nil {
		return fmt.Errorf("Failed to create database: %w", err)
	}

	// create clients vk, youtube, facebook instagram, etc.
	vkClient := vk.NewClient(
		cfg.VK.Host,
		cfg.VK.BasePath,
		cfg.VK.Version,
		cfg.VK.Token,
	)

	a, err := app.NewApp(&cfg, storage, vkClient)
	if err != nil {
		return fmt.Errorf("Failed to create app: %w", err)
	}
	a.Run()

	shutdown.Graceful(
		[]os.Signal{os.Interrupt, syscall.SIGTERM},
		storage,
	)
	return nil
}
