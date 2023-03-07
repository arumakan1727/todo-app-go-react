package pgsql

import (
	"context"
	"testing"
	"time"

	"github.com/arumakan1727/todo-app-go-react/clock"
	"github.com/arumakan1727/todo-app-go-react/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreUser(t *testing.T) {
	ctx := context.Background()
	clk := clock.GetFixedClocker()
	r := newRepositoryForTest(t, ctx, clk)

	require.NoError(t, r.TruncateAll(ctx))

	testcase := []struct {
		name     string
		in       domain.User
		checkErr func(got error)
	}{
		{
			name: "success",
			in: domain.User{
				ID:          0,
				Role:        "Role-1",
				Email:       "test@example.com",
				PasswdHash:  []byte("PasswdHash-1"),
				DisplayName: "DisplayName-1",
				CreatedAt:   time.Time{},
			},
			checkErr: func(got error) {
				assert.NoError(t, got, "testcase-1")
			},
		},
		{
			name: "fail with already registered email (case insentive)",
			in: domain.User{
				ID:          0,
				Role:        "Role-2",
				Email:       "tEsT@example.com",
				PasswdHash:  []byte("PasswdHash-2"),
				DisplayName: "DisplayName-2",
				CreatedAt:   time.Time{},
			},
			checkErr: func(got error) {
				assert.ErrorIs(t, got, domain.ErrAlreadyExits, "testcase-2")
				assert.ErrorContains(t, got, "email", "testcase-2")
			},
		},
		{
			name: "success with yet unregistered email",
			in: domain.User{
				ID:          1727,
				Role:        "Role-3",
				Email:       "foo@example.com",
				PasswdHash:  []byte("PasswdHash-3"),
				DisplayName: "DisplayName-3",
				CreatedAt:   time.Time{},
			},
			checkErr: func(got error) {
				assert.NoError(t, got, "testcase-3")
			},
		},
	}

	for i, tt := range testcase {
		t.Logf("testcase-%02d: [Should %s]", i+1, tt.name)
		var got domain.User = tt.in // copy

		err := r.StoreUser(ctx, &got)
		tt.checkErr(err)
		if err == nil {
			assert.NotZero(t, got.ID)
			assert.Equal(t, tt.in.Role, got.Role)
			assert.Equal(t, tt.in.Email, got.Email)
			assert.Equal(t, tt.in.PasswdHash, got.PasswdHash)
			assert.Equal(t, tt.in.DisplayName, got.DisplayName)
			assert.Equal(t, clk.Now(), got.CreatedAt)
			assert.Equal(t, time.UTC, got.CreatedAt.Location())
		}
	}
}

func clearAndInsertUsers(t *testing.T, ctx context.Context, r domain.Repository) []domain.User {
	t.Helper()
	require.NoError(t, r.TruncateAll(ctx))
	us := []domain.User{
		{
			Role:        "admin",
			Email:       "user0@example.com",
			PasswdHash:  []byte("PasswdHash-0"),
			DisplayName: "DisplayName-0",
		},
		{
			Role:        "user",
			Email:       "user1@example.com",
			PasswdHash:  []byte("PasswdHash-1"),
			DisplayName: "DisplayName-1",
		},
		{
			Role:        "user",
			Email:       "user2@example.com",
			PasswdHash:  []byte("PasswdHash-2"),
			DisplayName: "DisplayName-2",
		},
	}
	for i := range us {
		err := r.StoreUser(ctx, &us[i])
		if err != nil {
			t.Fatalf("clearAndInsertUsers: failed to StoreUser: i=%d, %#v", i, err)
		}
	}
	return us
}

func TestReadUsers(t *testing.T) {
	ctx := context.Background()
	clk := clock.GetFixedClocker()
	r := newRepositoryForTest(t, ctx, clk)

	users := clearAndInsertUsers(t, ctx, r)

	t.Run("ListUsers-OK (should not contain PasswdHash)", func(t *testing.T) {
		t.Parallel()
		got, err := r.ListUsers(ctx)
		if assert.NoError(t, err) {
			// 返されるパスワードハッシュはnilでなければならない
			wants := make([]domain.User, len(users))
			copy(wants, users)
			for i := range wants {
				wants[i].PasswdHash = nil
			}
			assert.Equal(t, wants, got)
		}
	})
	t.Run("GetUserByEmail-OK", func(t *testing.T) {
		t.Parallel()
		got, err := r.GetUserByEmail(ctx, "user0@example.com")
		if assert.NoError(t, err) {
			assert.Equal(t, users[0], got)
			assert.Equal(t, time.UTC, got.CreatedAt.Location())
		}
	})
	t.Run("GetUserByEmail-NotFound", func(t *testing.T) {
		t.Parallel()
		_, err := r.GetUserByEmail(ctx, "hoge@example.com")
		assert.ErrorIs(t, err, domain.ErrNotFound)
	})
}
