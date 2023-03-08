package domain

import "time"

type Role = string

type AuthMaterial struct {
	UID  UserID
	Role Role
}

type AuthUsecase interface {
	GetAuthTokenMaxAge() time.Duration
	IssueAuthToken(ctx Ctx, email, passwd string) (AuthToken, error)
	ValidateAuthToken(Ctx, AuthToken) (AuthMaterial, error)
	RevokeAuthToken(Ctx, AuthToken) error
}

func (am *AuthMaterial) IsAdmin() bool {
	return am.Role == "admin"
}
