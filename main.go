package main

import (
	"database/sql"
	"github.com/go-co-op/gocron"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"os"
	"time"

	"LeavingEmployees/bot"
)

func main() {
	log.SetReportCaller(true)
	log.SetFormatter(&log.JSONFormatter{})

	err := godotenv.Load(".env")
	if err != nil {
		log.Error("Error loading .env file")
		return
	}

	dbConnectionString := os.Getenv("DB_CONNECTION_STRING")
	mysqlConn, err := sql.Open("mysql", dbConnectionString)
	if err != nil {
		log.Error(err)
		return
	}

	defer func() {
		err := mysqlConn.Close()
		if err != nil {
			log.Error(err)
			return
		}
	}()

	token := os.Getenv("TOKEN")
	api := slack.New(token)

	channelID := os.Getenv("CHANNEL_ID")

	s := gocron.NewScheduler(time.UTC)
	_, err = s.Cron("0 0 * * *").Do(bot.FindLeavingEmployees, mysqlConn, api, channelID)
	if err != nil {
		log.Error(err)
		return
	}
	s.StartBlocking()
}
