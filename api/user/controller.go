package user

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"LeavingEmployees/database/models"
)

type userController struct {
	dbConnection *sql.DB
}

func NewUserController(conn *sql.DB) IUser {
	return &userController{
		dbConnection: conn,
	}
}

func (a *userController) GetUsers(c *gin.Context) {
	users, err := models.SelectAll(a.dbConnection)
	if err != nil {
		log.Error(err)
	}

	c.IndentedJSON(http.StatusOK, users)
}
