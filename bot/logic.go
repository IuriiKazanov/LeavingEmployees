package bot

import (
	"database/sql"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"

	"LeavingEmployees/database/models"
)

func FindLeavingEmployees(dbConnection *sql.DB, api *slack.Client, channelID string) error {
	usersDB, err := models.SelectAll(dbConnection)
	if err != nil {
		log.Error(err)
		return err
	}

	usersSlack, err := api.GetUsers()
	if err != nil {
		log.Error(err)
		return err
	}

	var leavingUsers []string
	for _, userSlack := range usersSlack {
		if !userSlack.IsAppUser {
			continue
		}
		isUserContains := false
		for _, userDB := range usersDB {
			if userSlack.ID == userDB.ID {
				isUserContains = true
				if userSlack.Deleted && !userDB.IsDeleted {
					leavingUsers = append(leavingUsers, userDB.Name)
					err := models.UpdateStatus(dbConnection, userDB)
					if err != nil {
						log.Error(err)
					}
					break
				}
			}
		}
		if !isUserContains {
			user := models.User{
				ID:          userSlack.ID,
				WorkspaceID: userSlack.TeamID,
				IsDeleted:   userSlack.Deleted,
				Name:        userSlack.Name,
			}
			err := models.Insert(dbConnection, user)
			if err != nil {
				log.Error(err)
			}
		}
	}

	var message string
	switch len(leavingUsers) {
	case 0:
		message = "Nobody quit!"
	case 1:
		message = fmt.Sprintf("%v quit!", leavingUsers[0])
	default:
		for _, user := range leavingUsers[:len(leavingUsers) - 1] {
			message += user + ", "
		}
		message += fmt.Sprintf("%v quit!", leavingUsers[len(leavingUsers) - 1])
	}

	err = SendMessage(api, channelID, message)
	if err != nil {
		log.Error(err)
	}

	return nil
}
