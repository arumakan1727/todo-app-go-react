package domain

import "time"

type KVS interface {
	SaveAuth(ctx Ctx, a AuthToken, uid UserID, expiration time.Duration) error
	FetchAuth(Ctx, AuthToken) (UserID, error)
	DeleteAuth(Ctx, AuthToken) error
}
