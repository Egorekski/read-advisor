package main

import (
	"flag"
	"log"
	"os"
	tgClient "read-advise-links-bot/internal/clients/telegram"
	eventconsumer "read-advise-links-bot/internal/consumer/event-consumer"
	telegram "read-advise-links-bot/internal/events/telegram"
	"read-advise-links-bot/internal/storage/files"
	"strconv"
)

func main() {

	eventsProcessor := telegram.New(
		tgClient.New(os.Getenv("TELEGRAM_BOT_HOST"), mustToken()),
		files.New(os.Getenv("STORAGE_PATH")),
	)

	log.Printf("server started")
	batchSize, _ := strconv.Atoi(os.Getenv("BATCH_SIZE"))
	consumer := eventconsumer.New(eventsProcessor, eventsProcessor, batchSize)
	if err := consumer.Start(); err != nil {
		log.Fatalf("failed to start consumer: %v", err)
	}
	//fetcher := fetcher.New(tgClient)

	// processor := processor.New(tgClient)

	// consumer.Start(fetcher, processor)

}

func mustToken() string {
	token := flag.String(
		"token-bot-token",
		"",
		"token for access telegram bot",
	)

	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}

	return *token
}
