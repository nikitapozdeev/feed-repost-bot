package app

import (
	"fmt"

	"github.com/nikitapozdeev/feed-repost-bot/internal/clients/vk"
	"github.com/nikitapozdeev/feed-repost-bot/internal/storage"
	tele "gopkg.in/telebot.v3"
)

// app implements App interface
type app struct {
	bot     *tele.Bot
	storage storage.Storage
	vk      *vk.Client
}

type App interface {
	Run()
}

// NewApp creates new application and setups handlers
func NewApp(bot *tele.Bot, storage storage.Storage, vk *vk.Client) (App, error) {
	app := app{
		bot:     bot,
		storage: storage,
		vk:      vk,
	}
	app.handlers()
	return &app, nil
}

// Run start the application
func (a *app) Run() {
	go func() {
		a.bot.Start()
	}()
	fmt.Println("Listening for clients...")
}
