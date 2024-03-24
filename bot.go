package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/rest"
)

func main() {
	log.Println("starting bot...")

	// Create client
	client, err := disgo.New(
		os.Getenv("BOT_TOKEN"),
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(gateway.IntentGuildMessages, gateway.IntentMessageContent),
		),
		bot.WithEventListenerFunc(onMessageCreate),
	)

	if err != nil {
		log.Fatalf("failed to start: %v", err)
	}

	defer client.Close(context.Background())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to gateway
	if err = client.OpenGateway(ctx); err != nil {
		log.Fatalf("failed to connect to gateway: %v", err)
	}

	log.Println("connected!")

	// Wait for interrupt
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
	<-s

	log.Println("disconnecting...")
}

func onMessageCreate(event *events.MessageCreate) {
	// Ignore all bots
	if event.Message.Author.Bot {
		return
	}

	// Check if message is poll
	if event.Message.Poll != nil {
		// Delete message
		event.Client().Rest().DeleteMessage(event.ChannelID, event.MessageID, rest.WithReason("Poll"))
	}
}
