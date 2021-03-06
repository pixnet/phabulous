package controllers

import (
        "regexp"
	"github.com/Sirupsen/logrus"
	"github.com/etcinit/gonduit/constants"
	"github.com/etcinit/gonduit/requests"
	"github.com/pixnet/phabulous/app/bot"
	"github.com/pixnet/phabulous/app/factories"
	"github.com/pixnet/phabulous/app/messages"
	"github.com/pixnet/phabulous/app/resolvers"
	"github.com/gin-gonic/gin"
	"github.com/jacobstr/confer"
)

// FeedController handles feed webhook routes
type FeedController struct {
	Config       *confer.Config                  `inject:""`
	Slacker      *bot.SlackService               `inject:""`
	Factory      *factories.GonduitFactory       `inject:""`
	Commits      *resolvers.CommitResolver       `inject:""`
	Tasks        *resolvers.TaskResolver         `inject:""`
	Differential *resolvers.DifferentialResolver `inject:""`
	Logger       *logrus.Logger                  `inject:""`
}

// Register registers the route handlers for this controller
func (f *FeedController) Register(r *gin.RouterGroup) {
	front := r.Group("/feed")
	{
		front.POST("/receive", f.postReceive)
	}
}

func (f *FeedController) postReceive(c *gin.Context) {
	conduit, err := f.Factory.Make()

	if err != nil {
		panic(err)
	}

	c.Request.ParseForm()

	res, err := conduit.PHIDQuerySingle(
		string(c.Request.PostForm.Get("storyData[objectPHID]")),
	)

	if err != nil {
		panic(err)
	}

        storyText := c.Request.PostForm.Get("storyText")

        re := regexp.MustCompile("(r([A-Z]+)([a-z0-9]{12})): (.+) (\\(.+\\))")
        storyText = re.ReplaceAllString(storyText, "$1: `$4`")

	if res.URI != "" {
                re = regexp.MustCompile("(r([A-Z]+)([a-z0-9]{12})):")
                storyText = re.ReplaceAllString(storyText, "<" + res.URI + "|$1>:")
	}


        phid := string(c.Request.PostForm.Get("storyData[objectPHID]"))
	commits, err := conduit.DiffusionQueryCommits(
		requests.DiffusionQueryCommitsRequest{
			PHIDs: []string{phid},
		},
	)

        r, _ := regexp.Compile("(rTEST([a-z0-9]{12}))")
        r2, _ := regexp.Compile("(added a comment|raised a concern|added inline comments|accepted)")
        if err == nil && ! r.MatchString(storyText) && r2.MatchString(storyText) {
                commit := commits.Data[phid]
                regEmail, _ := regexp.Compile("(([a-z_0-9][-._a-z0-9]*[a-z_0-9])@[a-z_0-9][-._a-z0-9]*[a-z_0-9].[a-z_0-9]{2,3})")
                commitAuthorName := regEmail.FindStringSubmatch(commit.AuthorEmail)[2]
                storyAuthorRes, err := conduit.PHIDQuerySingle(
                    string(c.Request.PostForm.Get("storyAuthorPHID")),
                )


                // Mention commit author
                slackuser, err := f.Slacker.GetUserByEmail(commit.AuthorEmail)
                if err == nil && storyAuthorRes.Name != commitAuthorName {
                    storyText += " <@" + string(slackuser.ID) + "|" + commitAuthorName + ">"
                }
        }

        iconEmoji := ""
        regComment, _ := regexp.Compile("(added a comment|added inline comments)")
        if regComment.MatchString(storyText) {
            iconEmoji = ":memo:"
        }

        regConcern, _ := regexp.Compile("raised a concern")
        if regConcern.MatchString(storyText) {
            iconEmoji = ":raising_hand:"
        }

        regAccepted, _ := regexp.Compile("accepted")
        if regAccepted.MatchString(storyText) {
            iconEmoji = ":white_check_mark:"
        }
        storyText += " " + iconEmoji

	phidType := constants.PhidType(res.Type)
	icon := messages.PhidTypeToIcon(phidType)

        icon = messages.StoryTextToIcon(storyText)

	f.Slacker.FeedPost(storyText, icon)

	switch phidType {
	case constants.PhidTypeCommit:
		channelName, err := f.Commits.Resolve(res.Name)
		if err != nil {
			f.Logger.Error(err)
		}

		if channelName != "" {
			f.Slacker.SimplePost(channelName, storyText, icon, false, iconEmoji)
		}
		break
	case constants.PhidTypeTask:
		channelName, err := f.Tasks.Resolve(res.PHID)
		if err != nil {
			f.Logger.Error(err)
		}

		if channelName != "" {
			f.Slacker.SimplePost(channelName, storyText, icon, false, iconEmoji)
		}
		break
	case constants.PhidTypeDifferentialRevision:
		channelName, err := f.Differential.Resolve(res.PHID)
		if err != nil {
			f.Logger.Error(err)
		}

		if channelName != "" {
			f.Slacker.SimplePost(channelName, storyText, icon, false, iconEmoji)
		}
		break
	}

	c.JSON(200, gin.H{
		"status": "success",
		"messages": []string{
			"OK",
		},
	})
}
