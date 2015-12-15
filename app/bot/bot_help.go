package bot

import (
	"github.com/pixnet/phabulous/app/messages"
	"github.com/nlopes/slack"
)

// HandleUsage shows usage tip.
func (b *Bot) HandleUsage(ev *slack.MessageEvent, matches []string) {
	b.Slacker.SimplePost(
		ev.Channel,
		"Hi. For usage information, type `help`.",
		messages.IconTasks,
		true,
	)
}

// HandleHelp shows help.
func (b *Bot) HandleHelp(ev *slack.MessageEvent, matches []string) {
	b.Slacker.SimplePost(
		ev.Channel,
		`Available commands:
    *summon Dxxx* (channel): Asks reviewers of a revision to review it.
    *lookup Txxx* (channel, DM): Looks up a task by its number.
    *lookup Dxxx* (channel, DM): Looks up a revision by its number.
    *help* (channel, DM): Shows this help.`,
		messages.IconTasks,
		true,
	)
}
