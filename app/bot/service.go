package bot

import (
	"errors"

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
) {
	user := s.Config.GetString("slack.username")

	if s.Bot != nil {
		user = s.Bot.slackInfo.User.Name
	}

	s.Slack.PostMessage(
		channelName,
		storyText,
		slack.PostMessageParameters{
			Username: user,
			IconURL:  string(icon),
			AsUser:   asUser,
                        Parse: "none",
		},
	)
}

// FeedPost posts a message to Slack on the default bot channel.
func (s *SlackService) FeedPost(storyText string) error {
	if s.GetFeedChannel() == "" {
		return ErrMissingFeedChannel
	}

	s.SimplePost(s.GetFeedChannel(), storyText, messages.IconDefault, false)
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
