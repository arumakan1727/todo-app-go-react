package domain

import (
	"time"
)

type UserID uint64
type AuthToken string

type UserUsecase interface {
	Store(ctx Ctx, email, passwd, displayName, role string) (User, error)
	List(Ctx) ([]User, error)
}

func (u *User) ApplyTimezone(loc *time.Location) {
	u.CreatedAt = u.CreatedAt.In(loc)
}
