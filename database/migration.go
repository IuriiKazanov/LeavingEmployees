package main

import (
	"database/sql"
	"flag"
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
}

func dbConnect() (*sql.DB, error) {
	dbConnectionString := os.Getenv("DB_CONNECTION_STRING")
	mysqlConn, err := sql.Open("mysql", dbConnectionString)
	if err != nil {
		return nil, err
	}

	return mysqlConn, nil
}
