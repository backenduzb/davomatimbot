package response

import "time"

type ProfileResponse struct {
	ID int `json:"id"`
	Username string `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	IsOnline bool `json:"is_online"`
	IsBanned bool `json:"is_banned"`
	IsAdmin bool `json:"is_admin"`
	LastSeen time.Time `json:"last_seen"`
	TelegramId string `json:"telegram_id"`
}