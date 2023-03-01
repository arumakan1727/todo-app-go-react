package pgsql

import (
	"context"
	"testing"
	"time"

	"github.com/arumakan1727/todo-app-go-react/clock"
	"github.com/arumakan1727/todo-app-go-react/domain"
	"github.com/stretchr/testify/assert"
)

func TestStoreUser(t *testing.T) {
	ctx := context.Background()
	clk := clock.GetFixedClocker(nil)
	r := newRepositoryForTest(t, ctx, clk)

	repoImpl := r.(*repository)
	clearTable(t, ctx, repoImpl, "users")

	testcase := []struct {
		in       domain.User
		checkErr func(got error)
	}{
		{
			in: domain.User{
				ID:          0,
				Role:        "user",
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
			in: domain.User{
				ID:          0,
				Role:        "user",
				Email:       "test@example.com",
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
			in: domain.User{
				ID:          1727,
				Role:        "user",
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
		t.Logf("testcase #%02d", i+1)
		var got domain.User = tt.in // copy

		err := r.StoreUser(ctx, &got)
		tt.checkErr(err)
		if err == nil {
			assert.NotZero(t, got.ID)
			assert.Equal(t, tt.in.Role, got.Role)
			assert.Equal(t, tt.in.Email, got.Email)
			assert.Equal(t, tt.in.PasswdHash, got.PasswdHash)
			assert.Equal(t, tt.in.DisplayName, got.DisplayName)
			assert.Equal(t, tt.in.CreatedAt, got.CreatedAt)
		}
	}
}

func prepareUserRecord(t *testing.T) {
	t.Helper()
}

func TestListUsers(t *testing.T) {

}

func TestGetUserByEmail(t *testing.T) {

}
