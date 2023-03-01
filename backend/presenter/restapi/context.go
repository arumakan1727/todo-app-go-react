package restapi

import (
	"context"
	"errors"
	"fmt"
)

type ctxKeyUserID struct{}

var ErrCtxGetValue = errors.New("failed to get value from ctx")

func newCtxWithUserID(ctx context.Context, uid UserID) context.Context {
	return context.WithValue(ctx, ctxKeyUserID{}, uid)
}

func getUserIDFromCtx(ctx context.Context) (UserID, error) {
	v, ok := ctx.Value(ctxKeyUserID{}).(UserID)
	if !ok {
		return 0, fmt.Errorf("UserID: %w", ErrCtxGetValue)
	}
	return v, nil
}
