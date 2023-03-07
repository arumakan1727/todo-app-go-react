package restapitest

import (
	"net/http"
	"testing"

	"github.com/arumakan1727/todo-app-go-react/domain"
	"github.com/arumakan1727/todo-app-go-react/presenter/restapi"
	oapi "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthToken(t *testing.T) {
	t.Parallel()

	var token domain.AuthToken

	t.Run("Issue-OK", func(t *testing.T) {
		resp, err := gClient.IssueAuthToken(gCtx, restapi.ReqCreateAuthToken{
			Email:    oapi.Email(gUser.Email),
			Password: gUser.Passwd,
		})
		require.NoError(t, err)
		defer closeBody(resp)

		assert.Equal(t, 200, resp.StatusCode)
		assert.Empty(t, readAll(t, resp.Body))

		cookies := resp.Cookies()
		require.Len(t, cookies, 1)

		c := cookies[0]
		assert.Equal(t, restapi.CookieKeyAuthToken, c.Name)
		assert.NotEmpty(t, c.Value)
		assert.Equal(t, "/", c.Path)
		assert.Equal(t, "", c.Domain)
		assert.True(t, c.HttpOnly)
		assert.Equal(t, http.SameSiteStrictMode, c.SameSite)
		assert.Greater(t, c.MaxAge, 0)

		token = domain.AuthToken(c.Value)
	})
	t.Run("Issue-Incorrect email", func(t *testing.T) {
		resp, err := gClient.IssueAuthToken(gCtx, restapi.ReqCreateAuthToken{
			Email:    "wrong@example.com",
			Password: gUser.Passwd,
		})
		require.NoError(t, err)
		defer closeBody(resp)

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		got := decodeRespAsSimpleErrorJSON(t, resp)
		assert.Equal(t, "incorrect email or password", got.Message)
	})
	t.Run("Issue-Incorrect password", func(t *testing.T) {
		resp, err := gClient.IssueAuthToken(gCtx, restapi.ReqCreateAuthToken{
			Email:    oapi.Email(gUser.Email),
			Password: "wrong-password",
		})
		require.NoError(t, err)
		defer closeBody(resp)

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		got := decodeRespAsSimpleErrorJSON(t, resp)
		assert.Equal(t, "incorrect email or password", got.Message)
	})
	t.Run("DeleteAuthToken-OK", func(t *testing.T) {
		resp, err := gClient.DeleteAuthToken(gCtx, addAuthTokenCookie(token))
		require.NoError(t, err)
		defer closeBody(resp)

		assert.Equal(t, 200, resp.StatusCode)
		assert.Empty(t, readAll(t, resp.Body))

		cookies := resp.Cookies()
		require.Len(t, cookies, 1)

		c := cookies[0]
		assert.Equal(t, restapi.CookieKeyAuthToken, c.Name)
		assert.Empty(t, c.Value)
		assert.LessOrEqual(t, c.MaxAge, 0)
	})
	t.Run("DeleteAuthToken-2nd", func(t *testing.T) {
		resp, err := gClient.DeleteAuthToken(gCtx, addAuthTokenCookie(token))
		require.NoError(t, err)
		defer closeBody(resp)

		// FIXME: すでに失効済みのトークンを DELETE しようとしたときの結果は401で良いのか？
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		got := decodeRespAsSimpleErrorJSON(t, resp)
		assert.Contains(t, got.Message, "Missing or invalid")
		assert.Contains(t, got.Message, "Token")
	})
}
