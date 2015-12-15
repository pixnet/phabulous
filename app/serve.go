package app

import (
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/pixnet/phabulous/app/bot"
	"github.com/jacobstr/confer"
)

// ServeService provides the serve command
type ServeService struct {
	Engine  *EngineService    `inject:""`
	Config  *confer.Config    `inject:""`
	Logger  *logrus.Logger    `inject:""`
	Slacker *bot.SlackService `inject:""`
}

// Run starts up the HTTP server
func (s *ServeService) Run(c *cli.Context) {
	s.Logger.Infoln("Starting up the server... (a.k.a. coffee time)")

	engine := s.Engine.New()

	go s.Slacker.BootRTM()

	// Figure out which port to use
	port := ":" + strconv.Itoa(s.Config.GetInt("server.port"))

	engine.Run(port)

	s.Logger.Infoln("✔︎ Done!")
}
