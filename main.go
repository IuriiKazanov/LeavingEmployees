package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jasonlvhit/gocron"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"

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

	session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-2"),
		Credentials: credentials.NewSharedCredentials(".aws.env", "default"),
	})

	svc := sqs.New(sess)

	var sqsUrl *string
	sqsUrl = flag.String("q","https://sqs.us-east-2.amazonaws.com/905195353451/iurii_kazanov-local.fifo", "")

	_, err = svc.SendMessage(&sqs.SendMessageInput{
		MessageGroupId: aws.String("123"),
		DelaySeconds: aws.Int64(0),
		MessageBody: aws.String("Information about current NY Times fiction bestseller for week of 12/11/2016."),
		QueueUrl:    sqsUrl,
	})
	if err != nil {
		log.Error(err)
		return
	}

	msgResult, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl: sqsUrl,
		MaxNumberOfMessages: aws.Int64(1),
		VisibilityTimeout:   aws.Int64(0),
	})

	fmt.Println(msgResult.String())

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

	s := gocron.NewScheduler()
	err = s.Every(1).Minute().Do(bot.FindLeavingEmployees, mysqlConn, api, channelID)
	if err != nil {
		log.Error(err)
		return
	}

	<-s.Start()
}
