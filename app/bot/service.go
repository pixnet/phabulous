package bot

import (
	"errors"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/pixnet/phabulous/app/factories"
	"github.com/pixnet/phabulous/app/messages"
	"github.com/jacobstr/confer"
	"github.com/nlopes/slack"
)

var (
	// ErrMissingFeedChannel is used when the feed channel is not configured.
	ErrMissingFeedChannel = errors.New("Missing feed channel")
)

// SlackService provides access to the Slack service.
type SlackService struct {
	Config  *confer.Config            `inject:""`
	Logger  *logrus.Logger            `inject:""`
	Factory *factories.GonduitFactory `inject:""`

	Slack *slack.Client
	Bot   *Bot
}

// SimplePost posts a simple message to Slack. Most parameters are set to
// defaults.
func (s *SlackService) SimplePost(
	channelName string,
	storyText string,
	icon messages.Icon,
	asUser bool,
        iconEmoji ...string,
) {
	user := s.Config.GetString("slack.username")

	if s.Bot != nil {
		user = s.Bot.slackInfo.User.Name
	}

        emoji := ""
        fmt.Println("iconEmoji count = ", len(iconEmoji))
        fmt.Println(iconEmoji)

        if len(iconEmoji) > 0 {
            emoji = iconEmoji[0]
            fmt.Println("emoji = " + emoji)
        }

	s.Slack.PostMessage(
		channelName,
		storyText,
		slack.PostMessageParameters{
			Username: user,
			IconURL:  string(icon),
                        IconEmoji: string(emoji),
			AsUser:   asUser,
                        Parse: "none",
		},
	)
}

// FeedPost posts a message to Slack on the default bot channel.
func (s *SlackService) FeedPost(storyText string, iconArgs ...messages.Icon) error {
	if s.GetFeedChannel() == "" {
		return ErrMissingFeedChannel
	}

        if len(iconArgs) > 0 {
                fmt.Println("len(iconArgs) > 0, ", iconArgs[0])
                s.SimplePost(s.GetFeedChannel(), storyText, iconArgs[0], false)
        } else {
                s.SimplePost(s.GetFeedChannel(), storyText, messages.IconDefault, false)
        }
	return nil
}

// GetFeedChannel returns the default channel for the bot.
func (s *SlackService) GetFeedChannel() string {
	return s.Config.GetString("channels.feed")
}

// GetUser return user given email
func (s *SlackService) GetUserByEmail(email string) (*slack.User, error) {
    users, err := s.Slack.GetUsers();
    if err != nil {
        return nil, err
    }

    for _, user := range users {
        if user.Profile.Email == email {
            return &user, nil
        }
    }

    return nil, nil
}
