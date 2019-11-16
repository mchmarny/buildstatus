package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/mchmarny/gcputil/env"
	"github.com/nlopes/slack"
)

var (
	sendStatus = []string{"WORKING", "SUCCESS", "FAILURE", "INTERNAL_ERROR", "TIMEOUT"}
	token      = env.MustGetEnvVar("SLACK_API_TOKEN", "")
	channel    = env.MustGetEnvVar("SLACK_BUILD_STATUS_CHANNEL", "")
)

func send(msg *CloudBuildNotification) error {

	if msg == nil {
		return fmt.Errorf("Null message on send: %v", msg)
	}

	// check if status is to be sent
	if !isStatusForSend(msg.Status) {
		log.Printf("Status not for send: %s != [%s]", msg.Status, strings.Join(sendStatus, ","))
		return nil
	}

	api := slack.New(token)

	a1 := slack.Attachment{
		Title:     "Trigger: ",
		TitleLink: msg.LogURL,
	}
	a1.Fields = []slack.AttachmentField{
		slack.AttachmentField{
			Title: "Tag",
			Value: fmt.Sprintf("Git repo *%s* was tagged: *%s*",
				msg.Source.RepoSource.RepoName,
				msg.Source.RepoSource.TagName),
		},
		slack.AttachmentField{
			Title: "Status",
			Value: msg.Status,
		},
	}

	_, _, err := api.PostMessage(channel,
		slack.MsgOptionText("Cloud Build Status", false),
		slack.MsgOptionAttachments(a1))

	return err

}

func isStatusForSend(a string) bool {
	for _, b := range sendStatus {
		if b == a {
			return true
		}
	}
	return false
}
