package main

import (
	"database/sql"
	"os"

	"github.com/jasonlvhit/gocron"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	slackSDK "github.com/slack-go/slack"

	"LeavingEmployees/slack"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Error("Error loading .env file")
		return
	}

	dbConnectionString := os.Getenv("DB_CONNECTION_STRING")
	mysqlConn, err := sql.Open("mysql", dbConnectionString)
	if err != nil {
		log.Error(err.Error())
		return
	}

	defer func() {
		err := mysqlConn.Close()
		if err != nil {
			log.Error(err.Error())
			return
		}
	}()

	token := os.Getenv("TOKEN")
	api := slackSDK.New(token)

	s := gocron.NewScheduler()
	err = s.Every(1).Minute().Do(slack.SendMessage, api, mysqlConn)
	//err = s.Every(1).Minute().Do(slack.SendSlackNotification, "https://hooks.slack.com/services/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")
	if err != nil {
		log.Error(err.Error())
		return
	}

	<-s.Start()
}
