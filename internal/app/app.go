package app

import (
	"fmt"
	"log"
	"os"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nikitapozdeev/feed-repost-bot/internal/config"
	"github.com/nikitapozdeev/feed-repost-bot/internal/db"
	"github.com/nikitapozdeev/feed-repost-bot/internal/model"
	"github.com/nikitapozdeev/feed-repost-bot/pkg/shutdown"
	tele "gopkg.in/telebot.v3"
)

type app struct {
	cfg *config.Config
	db  *db.DB
}

type App interface {
	Run()
}

var globalApp *app

func NewApp(cfg *config.Config, db *db.DB) (App, error) {
	return &app{
		cfg: cfg,
		db:  db,
	}, nil
}

func StartHandler(c tele.Context) error {
	return c.Send("TODO: send usage")
}

func PingHandler(c tele.Context) error {
	return c.Send("pong")
}

func AddHandler(c tele.Context) error {
	subscription := model.Subscription{
		ClientID: c.Sender().ID,
		FeedLink: c.Args()[0],
	}

	if err := globalApp.db.Add(subscription); err != nil {
		fmt.Println("[ERROR]: adding subscription failed: %w", err)
		c.Send("fail")
	}

	return c.Send("added")
}

func RemoveHandler(c tele.Context) error {
	return c.Send("remove")
}

func (a *app) Run() {
	globalApp = a
	tgSettings := tele.Settings{
		Token:  a.cfg.Token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := tele.NewBot(tgSettings)
	if err != nil {
		log.Fatal(err)
		return
	}

	// routes
	bot.Handle("/start", StartHandler)
	bot.Handle("/ping", PingHandler)
	bot.Handle("/add", AddHandler)
	bot.Handle("/remove", RemoveHandler)

	fmt.Println("Listening for clients...")
	
	go func() {
		bot.Start()
	}()

	shutdown.Graceful(
		[]os.Signal{os.Interrupt, syscall.SIGTERM}, 
		a.db,
	)
}