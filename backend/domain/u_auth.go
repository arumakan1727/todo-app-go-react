package domain

import "time"

type AuthUsecase interface {
	GetAuthTokenMaxAge() time.Duration
	IssueAuthToken(ctx Ctx, email, passwd string) (AuthToken, error)
	ValidateAuthToken(Ctx, AuthToken) (UserID, error)
	RevokeAuthToken(Ctx, AuthToken) error
}
