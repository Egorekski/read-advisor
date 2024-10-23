package event_consumer

import (
	"log"
	"read-advise-links-bot/internal/events"
	"time"
)

const (
	Second = 1 * time.Second
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

func New(fetcher events.Fetcher, processor events.Processor, batchSize int) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (c Consumer) Start() error {
	for {
		gotEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			log.Printf("[ERROR] consumer: %s", err.Error())

			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(Second)

			continue
		}

		if err := c.handleEvent(gotEvents); err != nil {
			log.Printf("[ERROR] consumer: %s", err.Error())

			continue
		}
	}
}

/*
TODO: solve problems:
1. events lost - solution: retry, возвращение в хранилище, fallback, confirmation for fetcher
2. lost all pack - solution: stop after first error
3. parallel processing - solution: use goroutines
*/

func (c Consumer) handleEvent(events []events.Event) error {
	//sync.WaitGroup{}
	for _, event := range events {
		log.Printf("[INFO] gor new event: %s", event.Text)

		if err := c.processor.Process(event); err != nil {
			log.Printf("[ERROR] can not handle event: %s", err.Error())

			continue
		}
	}
	return nil
}
