package restapitest

import (
	"net/http"
	"strings"
	"testing"

	"github.com/arumakan1727/todo-app-go-react/presenter/restapi"
	oapi "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	t.Parallel()

	t.Run("OK", func(t *testing.T) {
		req := restapi.ReqCreateUser{
			DisplayName: "DisplayName1",
			Email:       oapi.Email(uuid.New().String() + "@example.com"),
			Password:    "Password1",
		}
		resp, err := gClient.CreateUser(gCtx, req)
		defer closeBody(resp)
		require.NoError(t, err)

		assert.Equal(t, 200, resp.StatusCode)
		got := decodeRespAsJSON[restapi.RespUser](t, resp)
		want := restapi.RespUser{
			CreatedAt:   gClock.Now(),
			DisplayName: req.DisplayName,
			Email:       req.Email,
			Id:          got.Id,
		}
		assert.Equal(t, want, got)
		assert.NotZero(t, got.Id)
	})
	t.Run("Conflict with same email (case insentive)", func(t *testing.T) {
		p := strings.Split(gUser.Email, "@")
		email := strings.ToUpper(p[0]) + "@" + p[1]

		req := restapi.ReqCreateUser{
			DisplayName: "DisplayName2",
			Email:       oapi.Email(email),
			Password:    "Password2",
		}
		resp, err := gClient.CreateUser(gCtx, req)
		defer closeBody(resp)
		require.NoError(t, err)

		assert.Equal(t, http.StatusConflict, resp.StatusCode)
		got := decodeRespAsSimpleErrorJSON(t, resp)
		assert.Contains(t, got.Message, "email")
		assert.Contains(t, got.Message, req.Email)
	})
}

func TestListUsersForAdmin(t *testing.T) {
	t.Parallel()

	t.Run("OK", func(t *testing.T) {
		t.Parallel()
		resp, err := gClient.ListUsersForAdmin(gCtx, addAuthTokenCookie(gAdmin.Token))
		defer closeBody(resp)
		require.NoError(t, err)

		users, err := gRepo.ListUsers(gCtx)
		require.NoError(t, err)

		want := restapi.RespUserList{
			Items:      make([]restapi.RespUser, 0, len(users)),
			TotalCount: len(users),
		}
		for i := range users {
			u := &users[i]
			want.Items = append(want.Items, restapi.RespUser{
				CreatedAt:   u.CreatedAt,
				DisplayName: u.DisplayName,
				Email:       oapi.Email(u.Email),
				Id:          u.ID,
			})
		}

		assert.Equal(t, 200, resp.StatusCode)
		got := decodeRespAsJSON[restapi.RespUserList](t, resp)
		assert.Equal(t, want, got)
	})
	t.Run("NotFound with normal user", func(t *testing.T) {
		t.Parallel()
		resp, err := gClient.ListUsersForAdmin(gCtx, addAuthTokenCookie(gUser.Token))
		defer closeBody(resp)
		require.NoError(t, err)

		assert.Equal(t, 404, resp.StatusCode)
		got := decodeRespAsSimpleErrorJSON(t, resp)
		assert.Equal(t, "Not Found", got.Message)
	})
}
