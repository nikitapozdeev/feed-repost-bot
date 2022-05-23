package app

import (
	"fmt"

	"github.com/nikitapozdeev/feed-repost-bot/internal/poller"
	"github.com/nikitapozdeev/feed-repost-bot/internal/storage"
	tele "gopkg.in/telebot.v3"
)

// app implements App interface
type app struct {
	bot     *tele.Bot
	storage storage.Storage
	poller  *poller.Poller
}

type App interface {
	Run()
}

// NewApp creates new application and setups handlers
func NewApp(bot *tele.Bot, storage storage.Storage, poller *poller.Poller) (App, error) {
	app := app{
		bot:     bot,
		storage: storage,
		poller:  poller,
	}
	app.handlers()
	return &app, nil
}

// Run start the application
func (a *app) Run() {
	go func() {
		a.bot.Start()
	}()
	go func() {
		a.poller.Start()
	}()
	fmt.Println("Listening for clients...")
}
