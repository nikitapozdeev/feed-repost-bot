package main

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/nikitapozdeev/feed-repost-bot/internal/app"
	"github.com/nikitapozdeev/feed-repost-bot/internal/clients/vk"
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

	// create clients vk, youtube, facebook instagram, etc.
	vkClient := vk.NewClient(
		cfg.VK.Host, 
		cfg.VK.BasePath,
		cfg.VK.Version, 
		cfg.VK.Token,
	)

	a, err := app.NewApp(&cfg, db, vkClient)
	if err != nil {
		log.Fatal(err)
	}

	a.Run()
}