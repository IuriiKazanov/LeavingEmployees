package slack

import (
	"LeavingEmployees/database/models"
	"bytes"
	"database/sql"
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

func SendMessage(api *slack.Client, db *sql.DB) error {
	users, err := api.GetUsers()
	if err != nil {
		fmt.Printf("%s\n", err)
		return err
	}
	for _, user := range users {
		fmt.Printf("ID: %s, Name: %s, Deleted: %v\n", user.ID, user.Name, user.Deleted)
		user := models.User{
			ID:          user.ID,
			WorkspaceID: user.TeamID,
			IsDeleted:   user.Deleted,
			Name:        user.Name,
		}
		err := models.Insert(db, user)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
	}

	channelID, timestamp, err := api.PostMessage("U01U3T9324X", slack.MsgOptionText("Hello", false))

	if err != nil {
		fmt.Printf("%s\n", err)
		return err
	}
	fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
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
