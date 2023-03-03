package domain

import "time"

type UserID uint64
type AccessToken string

type UserUcase interface {
	Store(ctx Ctx, email, passwd, displayName string) (User, error)
	List(Ctx) ([]User, error)
	IssueAuthToken(ctx Ctx, email, passwd string) (AccessToken, error)
}

func (u *User) ApplyTimezone(loc *time.Location) {
	u.CreatedAt = u.CreatedAt.In(loc)
}
