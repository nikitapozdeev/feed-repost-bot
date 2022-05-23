package processor

import (
	"github.com/nikitapozdeev/feed-repost-bot/internal/producer"
	tele "gopkg.in/telebot.v3"
)

type Processor interface {
	Process(producer producer.Message, recipient int64) error
}

type MessageProcessor struct {
	bot *tele.Bot
}

func NewMessageProcessor(bot *tele.Bot) Processor {
	return &MessageProcessor{
		bot: bot,
	}
}

func (p *MessageProcessor) Process(message producer.Message, recipient int64) error {
	tgRecipient := tele.Chat{ID: recipient}
	htmlMessage, err := message.HTML()
	if err != nil {
		return err
	}

	msgOptions := tele.SendOptions{
		ParseMode:             tele.ModeHTML,
		DisableWebPagePreview: false,
	}
	p.bot.Send(&tgRecipient, htmlMessage, &msgOptions)

	return nil
}
