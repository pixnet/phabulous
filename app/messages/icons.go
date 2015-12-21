package messages

import (
        "regexp"
        "github.com/etcinit/gonduit/constants"
)

// Icon is the type of the PHID.
type Icon string

const (
	// IconDefault is used for regular messages.
	IconDefault Icon = "http://i.imgur.com/7Hzgo9Y.png"

	// IconCommits is used for commit-related messages.
	IconCommits Icon = "http://i.imgur.com/v8ReRKx.png"

	// IconTasks is used for task-related messages.
	IconTasks Icon = "http://i.imgur.com/jD7rf9x.png"

	// IconRevisions is used for revision-related messages.
	IconRevisions Icon = "http://i.imgur.com/NiPouYj.png"

        IconComment Icon = "https://slack.global.ssl.fastly.net/d4bf/img/emoji_2015_2/apple/1f4dd.png"

        IconAccepted Icon = "https://slack.global.ssl.fastly.net/d4bf/img/emoji_2015_2/apple/2714.png"

        IconConcern Icon = "https://slack.global.ssl.fastly.net/d4bf/img/emoji_2015_2/apple/1f64b.png"
)

// PhidTypeToIcon gets the matching icon for a PHID type.
func PhidTypeToIcon(phidType constants.PhidType) Icon {
	switch phidType {
	case constants.PhidTypeCommit:
		return IconCommits
	case constants.PhidTypeTask:
		return IconTasks
	case constants.PhidTypeDifferentialRevision:
		return IconRevisions
	default:
		return IconDefault
	}
}

// StoryTextToIcon gets the matching icon for a commit message.
func StoryTextToIcon(storyText string) Icon {
        regComment, _ := regexp.Compile("(added a comment|added inline comments)")
        if regComment.MatchString(storyText) {
            return IconComment
        }

        regConcern, _ := regexp.Compile("raised a concern")
        if regConcern.MatchString(storyText) {
            return IconConcern
        }

        regAccepted, _ := regexp.Compile("accepted")
        if regAccepted.MatchString(storyText) {
            return IconAccepted
        }

        return IconDefault
}
