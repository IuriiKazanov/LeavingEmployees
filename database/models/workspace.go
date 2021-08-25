package models

type Workspace struct {
	ID             string `json:"id"`
	IsActive       bool   `json:"isActive"`
	BotAccessToken string `json:"botAccessToken"`
}
