package mor

import (
	"fmt"

	"github.com/slack-go/slack"
	"github.com/vesicash/mor-api/external/request"
	"github.com/vesicash/mor-api/internal/config"
)

func SlackNotify(extReq request.ExternalRequest, channel, message string) error {
	if extReq.Test {
		return nil
	}

	api := slack.New(config.GetConfig().Slack.OauthToken)
	channelID := channel
	msg := message
	_, _, err := api.PostMessage(channelID, slack.MsgOptionText(msg, false))
	if err != nil {
		extReq.Logger.Error("error sending message to slack", err.Error())
		return fmt.Errorf("error sending message: %v", err)
	}

	return nil
}
