package domain

import "time"

type KVS interface {
	Close()
	SaveAuth(ctx Ctx, token AuthToken, r *AuthMaterial, expiration time.Duration) error
	FetchAuth(Ctx, AuthToken) (AuthMaterial, error)
	DeleteAuth(Ctx, AuthToken) error
}
