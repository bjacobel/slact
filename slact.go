package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bjacobel/slact/Godeps/_workspace/src/github.com/davecgh/go-spew/spew"
	"github.com/bjacobel/slact/Godeps/_workspace/src/github.com/joho/godotenv"
	"github.com/bjacobel/slact/Godeps/_workspace/src/github.com/nlopes/slack"
)

var spw = spew.NewDefaultConfig()

func main() {
	// Add some configuration to our JSON logger
	spw = &spew.ConfigState{Indent: "\t", MaxDepth: 5}

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	api := slack.New(os.Getenv("SLACK"))

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch msg.Data.(type) {

			case *slack.ReactionAddedEvent:
				go handleReaction(msg.Data.(*slack.ReactionAddedEvent), rtm)
			default:
				continue
			}
		}
	}
}

func handleReaction(reaction *slack.ReactionAddedEvent, rtm *slack.RTM) {
	rtm.SendMessage(
		rtm.NewOutgoingMessage(
			fmt.Sprintf("New reaction added :%s:", reaction.Reaction),
			reaction.Item.Item.Channel,
		),
	)
}
