package restapi

import (
	"context"
	"errors"
	"fmt"

	"github.com/arumakan1727/todo-app-go-react/domain"
	"github.com/labstack/echo/v4"
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

const echoCtxKeyAuthToken = "AuthToken"

func storeAuthTokenIntoCtx(ctx echo.Context, t domain.AuthToken) {
	ctx.Set(echoCtxKeyAuthToken, t)
}

func getAuthTokenFromCtx(ctx echo.Context) domain.AuthToken {
	return ctx.Get(echoCtxKeyAuthToken).(domain.AuthToken)
}
