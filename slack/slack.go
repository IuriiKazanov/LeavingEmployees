package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"io"
	"net/http"
)

type slackRequestBody struct {
	Text string `json:"text"`
}

func SendMessage(api *slack.Client, channelID, text string) error {
	channelID, timestamp, err := api.PostMessage(channelID, slack.MsgOptionText(text, false))
	if err != nil {
		fmt.Printf("%s\n", err)
		return err
	}

	fmt.Printf("Message successfully sent to channel %s at %s\n", channelID, timestamp)
	return nil
}

func SendSlackNotification(webhookUrl string) error {
	slackBody, err := json.Marshal(slackRequestBody{Text: "hello1"})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", webhookUrl, bytes.NewBuffer(slackBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error(err)
		}
	}(resp.Body)

	return nil
}
