package bot

import (
	"fmt"

	"github.com/pixnet/phabulous/app/messages"
	"github.com/nlopes/slack"
)

// HandleLookup looks up information about a Phabricator object.
func (b *Bot) HandleLookup(ev *slack.MessageEvent, matches []string) {
	b.slackRTM.SendMessage(b.slackRTM.NewTypingMessage(ev.Channel))

	conn, err := b.Slacker.Factory.Make()
	if err != nil {
		b.Excuse(ev, err)
		return
	}

	res, err := conn.PHIDLookupSingle(matches[1])
	if err != nil {
		b.Excuse(ev, err)
		return
	}

	b.Slacker.SimplePost(
		ev.Channel,
		fmt.Sprintf("*%s* (%s): %s", res.FullName, res.Status, res.URI),
		messages.IconTasks,
		true,
	)
}
