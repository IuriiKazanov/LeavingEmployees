package user

import (
	"github.com/gin-gonic/gin"
)

type IUser interface {
	GetUsers(c *gin.Context)
}
