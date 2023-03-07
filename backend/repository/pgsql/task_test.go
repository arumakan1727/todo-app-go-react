package pgsql

import (
	"context"
	"testing"

	"github.com/arumakan1727/todo-app-go-react/clock"
	"github.com/arumakan1727/todo-app-go-react/domain"
	"github.com/arumakan1727/todo-app-go-react/optional"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaskRepository(t *testing.T) {
	ctx := context.Background()
	clk := clock.GetFixedClocker()
	r := newRepositoryForTest(t, ctx, clk)

	users := clearAndInsertUsers(t, ctx, r)

	ts := make([]domain.Task, 3)

	t.Run("StoreTask", func(t *testing.T) {
		in := []struct {
			uid   domain.UserID
			title string
		}{
			{
				uid:   users[0].ID,
				title: "Title-0",
			},
			{
				uid:   users[0].ID,
				title: "Title-1",
			},
			{
				uid:   users[1].ID,
				title: "Title-2",
			},
		}
		for i := range in {
			in := &in[i]
			var err error
			ts[i], err = r.StoreTask(ctx, in.uid, in.title)

			require.NoError(t, err)

			got := &ts[i]
			assert.NotZero(t, got.ID)
			assert.Equal(t, in.uid, got.UserID)
			assert.Equal(t, in.title, got.Title)
			assert.False(t, got.Done)
			assert.Equal(t, clk.Now(), got.CreatedAt)
		}
	})
	t.Run("GetTask-OK", func(t *testing.T) {
		got, err := r.GetTask(ctx, ts[0].UserID, ts[0].ID)
		if assert.NoError(t, err) {
			assert.Equal(t, ts[0], got)
		}
	})
	t.Run("GetTask-NotFound", func(t *testing.T) {
		_, err := r.GetTask(ctx, users[2].ID, ts[0].ID)
		assert.ErrorIs(t, err, domain.ErrNotFound, "Non-owner cannot get other user's task")
		_, err = r.GetTask(ctx, ts[0].UserID, ts[0].ID+9999)
		assert.ErrorIs(t, err, domain.ErrNotFound, "Cannot get un-existing task")
	})
	t.Run("PatchTask", func(t *testing.T) {
		{
			ts[0].Title = "updatedTitle-0"
		}
		{
			ts[1].Done = true
		}
		{
			ts[2].Title = "updatedTitle-2"
			ts[2].Done = true
		}
		testcases := []*struct {
			name     string
			uid      domain.UserID
			tid      domain.TaskID
			patch    domain.TaskPatch
			want     *domain.Task
			checkErr func(err error)
		}{
			{
				name:     "success (update only title)",
				uid:      ts[0].UserID,
				tid:      ts[0].ID,
				patch:    domain.TaskPatch{Title: optional.Some("updatedTitle-0")},
				want:     &ts[0],
				checkErr: func(err error) { require.NoError(t, err) },
			},
			{
				name:     "fail with empty patch",
				uid:      ts[1].UserID,
				tid:      ts[1].ID,
				patch:    domain.TaskPatch{},
				want:     nil,
				checkErr: func(err error) { assert.ErrorIs(t, err, domain.ErrEmptyPatch) },
			},
			{
				name:     "fail with not found (non-owner)",
				uid:      ts[0].UserID,
				tid:      ts[2].ID,
				patch:    domain.TaskPatch{Title: optional.Some("updateTitle-NotFound")},
				want:     nil,
				checkErr: func(err error) { assert.ErrorIs(t, err, domain.ErrNotFound) },
			},
			{
				name:     "success (update only done)",
				uid:      ts[1].UserID,
				tid:      ts[1].ID,
				patch:    domain.TaskPatch{Done: optional.Some(true)},
				want:     &ts[1],
				checkErr: func(err error) { require.NoError(t, err) },
			},
			{
				name:     "success (update title and done)",
				uid:      ts[2].UserID,
				tid:      ts[2].ID,
				patch:    domain.TaskPatch{Title: optional.Some("updatedTitle-2"), Done: optional.Some(true)},
				want:     &ts[2],
				checkErr: func(err error) { require.NoError(t, err) },
			},
		}
		for _, tt := range testcases {
			tt := tt
			t.Run("Should "+tt.name, func(t *testing.T) {
				got, err := r.PatchTask(ctx, tt.uid, tt.tid, tt.patch)
				tt.checkErr(err)
				if err == nil {
					require.Equal(t, *tt.want, got)
				}
			})
		}
	})
	t.Run("ListTask-all", func(t *testing.T) {
		// ts の前から2つは users[0] が所有しているタスク
		want := ts[0:2]
		got, err := r.ListTasks(ctx, users[0].ID, domain.TaskListFilter{})
		if assert.NoError(t, err) {
			assert.Equal(t, want, got)
		}
	})
	t.Run("ListTask-empty", func(t *testing.T) {
		// users[2] はタスクを何も所有していない
		got, err := r.ListTasks(ctx, users[2].ID, domain.TaskListFilter{})
		if assert.NoError(t, err) {
			assert.Empty(t, got)
		}
	})
	t.Run("ListTask-where(done==false)", func(t *testing.T) {
		// users[0] が所有しているタスクの中で ts の最初の要素だけが done=false
		want := ts[0:1]
		f := domain.TaskListFilter{
			DoneEq: optional.Some(false),
		}
		got, err := r.ListTasks(ctx, users[0].ID, f)
		if assert.NoError(t, err) {
			assert.Equal(t, want, got)
		}
	})
	t.Run("ListTask-where(done==true)", func(t *testing.T) {
		// users[0] が所有しているタスクの中で ts の2番目の要素だけが done=true
		want := ts[1:2]
		f := domain.TaskListFilter{
			DoneEq: optional.Some(true),
		}
		got, err := r.ListTasks(ctx, users[0].ID, f)
		if assert.NoError(t, err) {
			assert.Equal(t, want, got)
		}
	})
	t.Run("DeleteTask-OK", func(t *testing.T) {
		uid := ts[0].UserID
		tid := ts[0].ID

		_, err := r.GetTask(ctx, uid, tid)
		require.NoError(t, err, "Delete前はタスクを取得できる")

		got, err := r.DeleteTask(ctx, uid, tid)
		if assert.NoError(t, err) {
			assert.Equal(t, ts[0], got)
		}

		_, err = r.GetTask(ctx, uid, tid)
		assert.ErrorIs(t, err, domain.ErrNotFound, "Delete後はタスクが見つからない")
	})
	t.Run("DeleteTask-NotFound", func(t *testing.T) {
		task := &ts[1]

		_, err := r.GetTask(ctx, task.UserID, task.ID)
		require.NoError(t, err, "削除対象のtaskは存在する")

		uid := users[2].ID
		require.NotEqual(t, uid, task.UserID, "id=%dはtaskの所有者ではない", uid)

		_, err = r.DeleteTask(ctx, uid, task.ID)
		assert.ErrorIs(t, err, domain.ErrNotFound, "非所有者はErrNotFoundによりタスクを削除できない")

		_, err = r.DeleteTask(ctx, uid, task.ID+9999)
		assert.ErrorIs(t, err, domain.ErrNotFound, "存在しないタスクはErrNotFoundにより削除できない")

		_, err = r.GetTask(ctx, task.UserID, task.ID)
		assert.NoError(t, err, "Delete後にErrNotFoundが起きてもタスクは存在し、所有者は取得できる")
	})
}
