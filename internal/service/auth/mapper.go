package auth

import db "github.com/MuhAndriansyah/yorkshire/internal/database/sqlc"

func toUserOutput(u db.User) UserOutput {
	return UserOutput{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		AvatarURL: u.AvatarUrl,
		CreatedAt: u.CreatedAt.Time,
		UpdatedAt: u.UpdatedAt.Time,
	}
}
