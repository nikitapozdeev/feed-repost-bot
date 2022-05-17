package main

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/nikitapozdeev/feed-repost-bot/internal/app"
	"github.com/nikitapozdeev/feed-repost-bot/internal/config"
	"github.com/nikitapozdeev/feed-repost-bot/internal/db"
)

func main() {
	var cfg config.Config
	err := cleanenv.ReadConfig("./config.yml", &cfg)
	if err != nil {
		log.Fatal(err)
	}

  db, err := db.NewDB("store")
	if err != nil {
		log.Fatal("[ERROR] creating db failed: ", err)
		return
	}
  defer db.Close()

	a, err := app.NewApp(&cfg, db)
	if err != nil {
		log.Fatal(err)
	}

	a.Run()
}