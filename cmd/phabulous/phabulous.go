package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/pixnet/phabulous/app"
	"github.com/facebookgo/inject"
	"github.com/jacobstr/confer"
)

func main() {
	// Seed rand.
	rand.Seed(time.Now().UTC().UnixNano())

	// Create the configuration
	// In this case, we will be using the environment and some safe defaults
	config := confer.NewConfig()
	config.ReadPaths("config/main.yml", "config/main.production.yml")
	config.AutomaticEnv()

	// Create the logger.
	logger := logrus.New()

	if config.GetBool("server.debug") {
		logger.Level = logrus.DebugLevel
	} else {
		logger.Level = logrus.WarnLevel
	}

	// Next, we setup the dependency graph
	// In this example, the graph won't have many nodes, but on more complex
	// applications it becomes more useful.
	var g inject.Graph
	var phabulous app.Phabulous
	g.Provide(
		&inject.Object{Value: config},
		&inject.Object{Value: &phabulous},
		&inject.Object{Value: logger},
	)
	if err := g.Populate(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Boot the upper layers of the app.
	phabulous.Boot()

	// Setup the command line application
	app := cli.NewApp()
	app.Name = "phabulous"
	app.Usage = "A Phabricator bot in Go"

	// Set version and authorship info
	app.Version = "2.1.0-alpha1"
	app.Author = "Eduardo Trujillo <ed@chromabits.com>"

	// Setup the default action. This action will be triggered when no
	// subcommand is provided as an argument
	app.Action = func(c *cli.Context) {
		fmt.Println("Usage: phabulous [global options] command [command options] [arguments...]")
	}

	app.Commands = []cli.Command{
		{
			Name:    "serve",
			Aliases: []string{"s", "server", "listen"},
			Usage:   "Start the API server",
			Action:  phabulous.Serve.Run,
		},
		{
			Name: "slack",
			Subcommands: []cli.Command{
				{
					Name:   "test",
					Usage:  "Test that the slackbot works",
					Action: phabulous.SlackWorkbench.SendTestMessage,
				},
				{
					Name:   "resolve.commit",
					Usage:  "Test that that a commit can correctly be resolved into a channel",
					Action: phabulous.SlackWorkbench.ResolveCommitChannel,
				},
				{
					Name:   "resolve.task",
					Usage:  "Test that that a task can correctly be resolved into a channel",
					Action: phabulous.SlackWorkbench.ResolveTaskChannel,
				},
				{
					Name:   "resolve.revision",
					Usage:  "Test that that a revision can correctly be resolved into a channel",
					Action: phabulous.SlackWorkbench.ResolveRevisionChannel,
				},
			},
		},
	}

	// Begin
	app.Run(os.Args)
}
