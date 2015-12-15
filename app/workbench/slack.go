package workbench

import (
	"github.com/codegangsta/cli"
	"github.com/pixnet/phabulous/app/bot"
	"github.com/pixnet/phabulous/app/factories"
	"github.com/pixnet/phabulous/app/resolvers"
	"github.com/jacobstr/confer"
	"github.com/nlopes/slack"
)

// SlackWorkbenchService provides tools for debugging and testing that the
// slack integration works.
type SlackWorkbenchService struct {
	Config       *confer.Config                  `inject:""`
	Slacker      *bot.SlackService               `inject:""`
	Factory      *factories.GonduitFactory       `inject:""`
	Commits      *resolvers.CommitResolver       `inject:""`
	Tasks        *resolvers.TaskResolver         `inject:""`
	Differential *resolvers.DifferentialResolver `inject:""`
}

// SendTestMessage sends a test message to the feeds channel.
func (s *SlackWorkbenchService) SendTestMessage(c *cli.Context) {
	_, _, err := s.Slacker.Slack.PostMessage(
		s.Config.GetString("channels.feed"),
		"This is a test message. Please ignore.",
		slack.PostMessageParameters{
			Username: s.Config.GetString("slack.username"),
		},
	)

	if err != nil {
		panic(err)
	}
}

// ResolveCommitChannel attempts to get which channel name should a commit be
// posted to.
func (s *SlackWorkbenchService) ResolveCommitChannel(c *cli.Context) {
	if res, _ := s.Commits.Resolve(c.Args().First()); res != "" {
		println("Target channel: " + res)

		return
	}

	println("No target channel found.")
}

// ResolveTaskChannel attempts to get which channel name should a commit be
// posted to.
func (s *SlackWorkbenchService) ResolveTaskChannel(c *cli.Context) {
	if res, _ := s.Tasks.Resolve(c.Args().First()); res != "" {
		println("Target channel: " + res)

		return
	}

	println("No target channel found.")
}

// ResolveRevisionChannel attempts to get which channel name should a revision
// be posted to.
func (s *SlackWorkbenchService) ResolveRevisionChannel(c *cli.Context) {
	if res, _ := s.Differential.Resolve(c.Args().First()); res != "" {
		println("Target channel: " + res)

		return
	}

	println("No target channel found.")
}
