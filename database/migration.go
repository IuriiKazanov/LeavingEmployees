package main

import (
	"LeavingEmployees/database/models"
	"database/sql"
	"flag"
	"github.com/slack-go/slack"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetReportCaller(true)
	log.SetFormatter(&log.JSONFormatter{})

	isDowngrade := flag.Bool("d", false, "a bool")
	flag.Parse()

	err := godotenv.Load(".env")
	if err != nil {
		log.Error("Error loading .env file")
		return
	}

	if *isDowngrade {
		downgrade()
	} else {
		upgrade()
	}
}

func upgrade() {
	mysqlConn, err := dbConnect()
	if err != nil {
		log.Error(err)
	}

	defer func() {
		err := mysqlConn.Close()
		if err != nil {
			log.Error(err)
			return
		}
	}()

	_, err = mysqlConn.Exec(
		"ALTER TABLE user ADD COLUMN imageUrl VARCHAR(255) NULL",
	)
	if err != nil {
		log.Error(err)
	}

	token := os.Getenv("TOKEN")
	api := slack.New(token)

	usersSlack, err := api.GetUsers()
	if err != nil {
		log.Error(err)
	}

	usersDB, err := models.SelectAll(mysqlConn)
	if err != nil {
		log.Error(err)
	}

	for _, userSlack := range usersSlack {
		for _, userDB := range usersDB {
			if userSlack.ID == userDB.ID {
				_, err = mysqlConn.Exec(
					"UPDATE user SET imageUrl=? WHERE userID=?",
					userSlack.Profile.Image192,
					userSlack.ID,
				)
				if err != nil {
					log.Error(err)
				}
			}
		}
	}

}

func downgrade() {
	mysqlConn, err := dbConnect()
	if err != nil {
		log.Error(err)
	}

	defer func() {
		err := mysqlConn.Close()
		if err != nil {
			log.Error(err)
			return
		}
	}()

	_, err = mysqlConn.Exec(
		"ALTER TABLE user DROP COLUMN imageUrl",
	)
	if err != nil {
		log.Error(err)
	}
}

func dbConnect() (*sql.DB, error) {
	dbConnectionString := os.Getenv("DB_CONNECTION_STRING")
	mysqlConn, err := sql.Open("mysql", dbConnectionString)
	if err != nil {
		return nil, err
	}

	return mysqlConn, nil
}
