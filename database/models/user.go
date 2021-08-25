package models

import (
	"database/sql"
)

type User struct {
	ID          string `json:"id"`
	WorkspaceID string `json:"workspaceID"`
	IsDeleted   bool   `json:"isDeleted"`
	Name        string `json:"name"`
}

func Insert(db *sql.DB, user User) error {
	_, err := db.Exec(
		"INSERT INTO user VALUES (?, ?, ?, ?)",
		user.ID,
		user.WorkspaceID,
		user.IsDeleted,
		user.Name,
	)

	return err
}

func UpdateStatus(db *sql.DB, user User) error {
	_, err := db.Exec(
		"UPDATE user SET isDeleted = not isDeleted WHERE userID = ? and workspaceID = ?",
		user.ID,
		user.WorkspaceID,
	)

	return err
}

func SelectAll(db *sql.DB) ([]User, error) {
	var users []User
	results, err := db.Query("SELECT userID, workspaceID, isDeleted, name FROM user")
	if err != nil {
		return nil, err
	}

	for results.Next() {
		var user User
		err = results.Scan(&user.ID, &user.WorkspaceID, &user.IsDeleted, &user.Name)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
