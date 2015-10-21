package main

import (
	"log"
	"net/http"
	"os"

	"github.com/bjacobel/slact/Godeps/_workspace/src/github.com/bjacobel/fresh/runner/runnerutils"
	"github.com/bjacobel/slact/Godeps/_workspace/src/github.com/davecgh/go-spew/spew"
	"github.com/bjacobel/slact/Godeps/_workspace/src/github.com/go-martini/martini"
	"github.com/bjacobel/slact/Godeps/_workspace/src/github.com/joho/godotenv"
	"github.com/bjacobel/slact/Godeps/_workspace/src/github.com/martini-contrib/binding"
	"github.com/bjacobel/slact/Godeps/_workspace/src/github.com/martini-contrib/gorelic"
	"github.com/bjacobel/slact/Godeps/_workspace/src/github.com/martini-contrib/render"
	"github.com/bjacobel/slact/models"
)

var spw = spew.NewDefaultConfig()

func runnerMiddleware(w http.ResponseWriter, r *http.Request) {
	if runnerutils.HasErrors() {
		runnerutils.RenderError(w)
	}
}

func main() {
	// Add some configuration to our JSON logger
	spw = &spew.ConfigState{Indent: "\t", MaxDepth: 5}

	app := martini.Classic()
	app.Use(render.Renderer())

	err := godotenv.Load()
	if err != nil {
		// If it's the dev environment, use Fresh as a reloader
		app.Use(runnerMiddleware)
	} else {
		gorelic.InitNewrelicAgent(os.Getenv("NRKEY"), "slact", false)
		app.Use(gorelic.Handler)
	}

	app.Group("/v1", func(r martini.Router) {
		r.Post("/messages", binding.Form(models.WebhookMsg{}), func(msg models.WebhookMsg, r render.Render) {
			log.Println(msg.Text)
			r.JSON(200, nil)
		})
	})

	app.Run()
}
