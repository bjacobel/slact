package main

import (
	"log"
	"net/http"
	"os"

	"github.com/bjacobel/slact/Godeps/_workspace/src/github.com/bjacobel/fresh/runner/runnerutils"
	"github.com/bjacobel/slact/Godeps/_workspace/src/github.com/codegangsta/martini"
	"github.com/bjacobel/slact/Godeps/_workspace/src/github.com/codegangsta/martini-contrib/render"
	"github.com/bjacobel/slact/Godeps/_workspace/src/github.com/davecgh/go-spew/spew"
	"github.com/bjacobel/slact/Godeps/_workspace/src/github.com/joho/godotenv"
	"github.com/bjacobel/slact/Godeps/_workspace/src/github.com/martini-contrib/gorelic"
)

var spw = spew.NewDefaultConfig()
var app = martini.Classic()

func init() {
	// Add some configuration to our JSON logger
	spw = &spew.ConfigState{Indent: "\t", MaxDepth: 5}

	app.Use(render.Renderer())

	if os.Getenv("DEV_RUNNER") == "1" {
		// If it's the dev environment, use Fresh as a reloader
		app.Use(runnerMiddleware)

		// then load .env into ENV
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	} else {
		// Set up New Relic
		gorelic.InitNewrelicAgent(os.Getenv("NRKEY"), "slact", false)
		app.Use(gorelic.Handler)
	}
}

func runnerMiddleware(w http.ResponseWriter, r *http.Request) {
	if runnerutils.HasErrors() {
		runnerutils.RenderError(w)
	}
}

func main() {
	app.Group("/v1", func(r martini.Router) {
		// r.Post(`/message`, binding.Json(models.Message{}), func(msg models.Message, r render.Render) {
		// }
	})

	app.Run()
}
