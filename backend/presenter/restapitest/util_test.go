package restapitest

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/arumakan1727/todo-app-go-react/domain"
	"github.com/arumakan1727/todo-app-go-react/presenter/restapi"
	"github.com/arumakan1727/todo-app-go-react/presenter/restapi/client"
	"github.com/stretchr/testify/require"
)

func readAll(t *testing.T, r io.ReadCloser) string {
	t.Helper()
	defer r.Close()
	bs, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("ReadAll: %#v", err)
	}
	return string(bs)
}

func decodeRespAsSimpleErrorJSON(t *testing.T, r *http.Response) restapi.RespSimpleError {
	return decodeRespAsJSON[restapi.RespSimpleError](t, r)
}

func decodeRespAsJSON[T any](t *testing.T, r *http.Response) T {
	t.Helper()
	defer r.Body.Close()

	ct := r.Header.Get("Content-Type")
	require.True(t, strings.HasPrefix(ct, "application/json"), "DecodeRespAsJSON: want prefix 'application/json', but got %q", ct)

	var dest T
	if err := json.NewDecoder(r.Body).Decode(&dest); err != nil {
		t.Fatalf("DecodeRespAsJSON: %#v", err)
	}
	return dest
}

func addAuthTokenCookie(token domain.AuthToken) client.RequestEditorFn {
	return func(ctx context.Context, req *http.Request) error {
		req.AddCookie(&http.Cookie{
			Name:  restapi.CookieKeyAuthToken,
			Value: string(token),
		})
		return nil
	}
}

func closeBody(r *http.Response) {
	_ = r.Body.Close()
}
