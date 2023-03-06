package restapi

import (
	"context"
	"errors"
	"fmt"

	"github.com/arumakan1727/todo-app-go-react/domain"
)

type ctxKeyAuthMaterial struct{}

var ErrCtxGetValue = errors.New("failed to get value from ctx")

func newCtxWithAuthMaterial(ctx context.Context, am domain.AuthMaterial) context.Context {
	return context.WithValue(ctx, ctxKeyAuthMaterial{}, am)
}

func getAuthMaterialFromCtx(ctx context.Context) (domain.AuthMaterial, error) {
	am, ok := ctx.Value(ctxKeyAuthMaterial{}).(domain.AuthMaterial)
	if !ok {
		return am, fmt.Errorf("getAuthMaterialFromCtx: %w", ErrCtxGetValue)
	}
	return am, nil
}
