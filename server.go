package main

import (
	"LeavingEmployees/api/user"
	"database/sql"
	"github.com/gin-gonic/gin"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
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

	router := gin.New()
	router.Use(CORSMiddleware())

	userController := user.NewUserController(mysqlConn)

	v1 := router.Group("/api/v1")
	{
		users := v1.Group("/users")
		{
			users.GET("", userController.GetUsers)
		}
	}

	err = router.Run("localhost:8080")
	if err != nil {
		log.Error(err)
		return
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "GET")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
