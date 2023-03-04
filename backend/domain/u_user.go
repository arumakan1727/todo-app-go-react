package domain

import "time"

type UserID uint64
type AuthToken string

type UserUcase interface {
	Store(ctx Ctx, email, passwd, displayName string) (User, error)
	List(Ctx) ([]User, error)
	GetAuthTokenLifeSpan() time.Duration
	IssueAuthToken(ctx Ctx, email, passwd string) (AuthToken, error)
	ValidateAuthToken(Ctx, AuthToken) (UserID, error)
	RevokeAuthToken(Ctx, AuthToken) error
}

func (u *User) ApplyTimezone(loc *time.Location) {
	u.CreatedAt = u.CreatedAt.In(loc)
}
