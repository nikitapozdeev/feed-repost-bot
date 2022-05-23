package poller

import (
	"log"
	"time"

	"github.com/nikitapozdeev/feed-repost-bot/internal/processor"
	"github.com/nikitapozdeev/feed-repost-bot/internal/producer"
	"github.com/nikitapozdeev/feed-repost-bot/internal/storage"
)

type Poller struct {
	Rate      time.Duration
	Storage   storage.Storage
	Producer  producer.Producer
	Processor processor.Processor
}

func NewPoller(rate time.Duration, storage storage.Storage, producer producer.Producer, processor processor.Processor) *Poller {
	return &Poller{
		Rate:      rate,
		Storage:   storage,
		Producer:  producer,
		Processor: processor,
	}
}

func (p *Poller) Start() {
	for {
		p.tick()
		time.Sleep(p.Rate)
	}
}

func (p *Poller) tick() {
	subscriptions, err := p.Storage.GetAllActive()
	if err != nil {
		log.Printf("Poller error: %s", err)
		return
	}

	if len(subscriptions) == 0 {
		return
	}

	for _, subscription := range subscriptions {
		messages, err := p.Producer.Posts(subscription.FeedLink, 0, 1)
		if err != nil {
			log.Printf("Poller error: %s", err)
			continue
		}

		if len(messages) > 0 {
			p.processMessages(messages, subscription.ClientID)
		}
	}
}

func (p *Poller) processMessages(messages []producer.Message, recipient int64) error {
	for _, message := range messages {
		if err := p.Processor.Process(message, recipient); err != nil {
			//log.Printf("Failed to process post %d", message.ID)
			log.Printf("Error: %s", err)
		}
	}

	return nil
}
