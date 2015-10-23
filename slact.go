package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bjacobel/slact/Godeps/_workspace/src/github.com/davecgh/go-spew/spew"
	"github.com/bjacobel/slact/Godeps/_workspace/src/github.com/go-martini/martini"
	"github.com/bjacobel/slact/Godeps/_workspace/src/github.com/joho/godotenv"
	"github.com/bjacobel/slact/Godeps/_workspace/src/github.com/martini-contrib/render"
	"github.com/bjacobel/slact/Godeps/_workspace/src/github.com/nlopes/slack"
)

var spw = spew.NewDefaultConfig()

func main() {
	// Add some configuration to our JSON logger
	spw = &spew.ConfigState{Indent: "\t", MaxDepth: 5}

	// load .env into os.Getenv
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Slack API client
	api := slack.New(os.Getenv("SLACK"))

	// Martini router and renderer
	app := martini.Classic()
	app.Use(render.Renderer())

	// Open a connection to Slack Real Time Messaging API
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	// expose Martini endpoints for our data
	app.Group("/", func(r martini.Router) {
		r.Get("reactions", func(r render.Render) {
			r.JSON(200, nil)
		})
	})

	app.Run()

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
