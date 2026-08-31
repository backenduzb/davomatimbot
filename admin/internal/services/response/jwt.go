package responses

import (
	"admin/internal/schemas/user/response"
)

func NewJWTTokenResponse(token string) response.LoginResponse {
	return response.LoginResponse{
		Token: token,
	}
}