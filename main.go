package main

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jasonlvhit/gocron"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"

	"LeavingEmployees/bot"
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
	api := slack.New(token)

	channelID := os.Getenv("CHANNEL_ID")

	s := gocron.NewScheduler()
	err = s.Every(1).Minute().Do(bot.FindLeavingEmployees, mysqlConn, api, channelID)
	if err != nil {
		log.Error(err.Error())
		return
	}

	<-s.Start()
}
