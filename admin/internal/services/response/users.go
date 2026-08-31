package responses

import (
	"admin/internal/schemas/user/response"
	"admin/internal/models"
)

func NewProfileResponse(user models.User) response.ProfileResponse {
	return response.ProfileResponse{
		ID: int(user.ID),
		Username: user.Username,
		CreatedAt: user.CreatedAt,
		IsOnline: user.IsOnline,
		IsBanned: user.IsBanned,
		IsAdmin: user.IsAdmin,
		TelegramId: user.TelegramId,
		LastSeen: user.LastSeen,
	}
} 
