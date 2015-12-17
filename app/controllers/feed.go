package controllers

import (
        "encoding/json"
        "math/rand"
        "regexp"
        "time"
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

        seed := rand.NewSource(time.Now().UnixNano() + 1)
        rnd := rand.New(seed)

        animals := []string{ ":cat:", ":dog:", ":mouse:", ":hamster:", ":rabbit:", ":wolf:", ":frog:", ":tiger:", ":koala:", ":bear:", ":pig:", ":cow:", ":boar:", ":monkey_face:", ":horse:", ":camel:", ":sheep:", ":elephant:", ":panda_face:", ":snake:", ":bird:", ":baby_chick:", ":chicken:", ":penguin:", ":turtle:", ":bug:", ":honeybee:", ":ant:", ":snail:", ":tropical_fish:", ":whale:", ":dolphin:", ":dragon:", ":rooster:", ":dragon_face:", ":crocodile:", ":poodle:", ":octopus:" }

        symbol := animals[rnd.Intn(len(animals))]

        storyText := symbol + " " + c.Request.PostForm.Get("storyText")

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
        if err == nil && r.MatchString(storyText) {
                commit := commits.Data[phid]
                storyText += " [DEBUG]: AuthorEmail = " + commit.AuthorEmail + ", Committer = " + commit.Committer

                storyJson, err := json.Marshal(string(c.Request.PostForm.Get("storyData")))
                if (err == nil) {
                    storyText += ", storyData = ```" + string(storyJson) + "```"
                }
        }


	phidType := constants.PhidType(res.Type)
	icon := messages.PhidTypeToIcon(phidType)

	f.Slacker.FeedPost(storyText)

	switch phidType {
	case constants.PhidTypeCommit:
		channelName, err := f.Commits.Resolve(res.Name)
		if err != nil {
			f.Logger.Error(err)
		}

		if channelName != "" {
			f.Slacker.SimplePost(channelName, storyText, icon, false)
		}
		break
	case constants.PhidTypeTask:
		channelName, err := f.Tasks.Resolve(res.PHID)
		if err != nil {
			f.Logger.Error(err)
		}

		if channelName != "" {
			f.Slacker.SimplePost(channelName, storyText, icon, false)
		}
		break
	case constants.PhidTypeDifferentialRevision:
		channelName, err := f.Differential.Resolve(res.PHID)
		if err != nil {
			f.Logger.Error(err)
		}

		if channelName != "" {
			f.Slacker.SimplePost(channelName, storyText, icon, false)
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
