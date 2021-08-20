package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jasonlvhit/gocron"
	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	token := os.Getenv("TOKEN")
	api := slack.New(token)

	s := gocron.NewScheduler()
	err = s.Every(1).Minute().Do(sendMessage, api)
	if err != nil {
		log.Fatalf(err.Error())
	}
	<-s.Start()
}

func sendMessage(api *slack.Client) error {
	users, err := api.GetUsers()
	if err != nil {
		fmt.Printf("%s\n", err)
		return err
	}
	for _, user := range users {
		fmt.Printf("ID: %s, Name: %s, Deleted: %v\n", user.ID, user.Name, user.Deleted)
	}

	channelID, timestamp, err := api.PostMessage("U01U3T9324X", slack.MsgOptionText("Hello", false))

	if err != nil {
		fmt.Printf("%s\n", err)
		return err
	}
	fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
	return nil
}
