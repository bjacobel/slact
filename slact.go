package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

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
		log.Println("No .env, running in production mode, hope that's what you wanted")
	}

	creds := ""

	if os.Getenv("MGUSER") != "" && os.Getenv("MGPW") != "" {
		creds = fmt.Sprintf("%s:%s@", os.Getenv("MGUSER"), os.Getenv("MGPW"))
	}

	// connect to a mongodb session - local or through Docker
	session, err := mgo.DialWithTimeout(
		fmt.Sprintf("mongodb://%s%s/%s", creds, os.Getenv("MGHOST"), os.Getenv("MGDB")),
		time.Second,
	)

	if err == nil {
		defer session.Close()
		session.SetMode(mgo.Monotonic, true)
	} else {
		log.Fatal(err)
	}

	db := session.Clone().DB(os.Getenv("MGDB")).C(os.Getenv("MGCOLL"))

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
			var result []bson.M

			pipe := db.Pipe(
				[]bson.M{
					{
						"$sort": bson.M{"eventtimestamp": 1},
					},
					{
						"$group": bson.M{
							"_id": "$reaction",
							"reactions": bson.M{
								"$push": "$$CURRENT",
							},
							"count": bson.M{"$sum": 1},
						},
					},
					{
						"$sort": bson.M{"count": -1},
					},
				},
			)

			iter := pipe.Iter()
			err := iter.All(&result)

			if err != nil {
				log.Println(err.Error())
				r.JSON(400, bson.M{"error": err.Error()})
			} else {
				r.JSON(200, result)
			}

		})
	})

	go app.Run()

	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch msg.Data.(type) {

			case *slack.ReactionAddedEvent:
				go insertReaction(msg.Data.(*slack.ReactionAddedEvent), db)
			case *slack.ReactionRemovedEvent:
				go deleteReaction(msg.Data.(*slack.ReactionRemovedEvent), db)
			default:
				continue
			}
		}
	}
}

func insertReaction(reax *slack.ReactionAddedEvent, db *mgo.Collection) {
	db.Insert(&reax)
	// log.Printf("Added reaction %s\n", reax.Reaction)
}

func deleteReaction(reax *slack.ReactionRemovedEvent, db *mgo.Collection) {
	err := db.Remove(bson.M{"user": reax.User, "item.item.timestamp": reax.Item.Item.Timestamp})

	if err != nil {
		return
	}

	log.Println("Could not find right object to delete. Not a huge deal.")
}
