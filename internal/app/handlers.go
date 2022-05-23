package app

import (
	"fmt"

	"github.com/nikitapozdeev/feed-repost-bot/internal/storage"
	tele "gopkg.in/telebot.v3"
)

// handlers set bot handlers
func (a *app) handlers() {
	a.bot.Handle("/start", a.handlerStart)
	a.bot.Handle("/ping", a.handlerPing)
	a.bot.Handle("/add", a.handlerAdd)
	a.bot.Handle("/remove", a.handlerRemove)
	a.bot.Handle(tele.OnText, a.handlerText)
}

// handlerStart handles bot /start command
func (a *app) handlerStart(c tele.Context) error {
	return c.Send("TODO: send usage")
}

// handlerPing handles bot /ping command
func (a *app) handlerPing(c tele.Context) error {
	return c.Send("pong")
}

// handlerAdd handles bot /add command
func (a *app) handlerAdd(c tele.Context) error {
	subscription := storage.Subscription{
		ClientID: c.Sender().ID,
		FeedLink: c.Args()[0],
	}

	fmt.Println("Received add message from", c.Sender().ID)

	if err := a.storage.Add(subscription); err != nil {
		fmt.Println("[ERROR]: adding subscription failed: %w", err)
		c.Send("fail")
	}

	c.Send("Subscription added")
	return nil
}

// handlerRemove handles bot /remove command
func (a *app) handlerRemove(c tele.Context) error {
	return c.Send("remove")
}

// handlerText handles text messages
func (a *app) handlerText(c tele.Context) error {
	if c.Message().IsForwarded() && c.Message().OriginalChat.Type == tele.ChatChannel {
		clientId := c.Message().Sender.ID
		chatId := c.Message().OriginalChat.ID
		subscriptions, err := a.storage.Get(clientId)
		if err != nil {
			return err
		}

		if len(subscriptions) == 0 {
			c.Send("You have not subscribed to any feed yet")
			return nil
		}

		for _, subscribtion := range subscriptions {
			// TODO: we need to store clientId and chatId in separate fields
			// one client can post feed to multiply channels/chats
			subscribtion.ClientID = chatId
			subscribtion.IsActive = true
			a.storage.Update(subscribtion)
		}
	} else {
		c.Send("unknown command, see usage")
	}

	return nil
}
