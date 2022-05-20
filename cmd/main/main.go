package main

import (
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/nikitapozdeev/feed-repost-bot/internal/app"
	"github.com/nikitapozdeev/feed-repost-bot/internal/clients/vk"
	"github.com/nikitapozdeev/feed-repost-bot/internal/config"
	"github.com/nikitapozdeev/feed-repost-bot/internal/storage/sqldb"
	"github.com/nikitapozdeev/feed-repost-bot/pkg/shutdown"
	tele "gopkg.in/telebot.v3"
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

	// create storage
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

	// create telebot
	tbSettings := tele.Settings{
		Token:  cfg.Telegram.Token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}
	bot, err := tele.NewBot(tbSettings)

	// create and run the app
	a, err := app.NewApp(bot, storage, vkClient)
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
