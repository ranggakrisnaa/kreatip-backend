package converter

import (
	"github.com/kreatip/kreatip-backend/internal/entity"
	"github.com/kreatip/kreatip-backend/internal/model"
)

func ToUserResponse(u *entity.User) *model.UserResponse {
	return &model.UserResponse{
		ID:            u.ID,
		Email:         u.Email,
		Role:          u.Role,
		EmailVerified: u.EmailVerified,
		TOTPEnabled:   u.TOTPEnabled,
	}
}
