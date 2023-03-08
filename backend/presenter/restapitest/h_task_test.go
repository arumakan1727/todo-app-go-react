package restapitest

import (
	"testing"

	"github.com/arumakan1727/todo-app-go-react/domain"
	"github.com/arumakan1727/todo-app-go-react/presenter/restapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func toRespTask(t *domain.Task) restapi.RespTask {
	return restapi.RespTask{
		CreatedAt: t.CreatedAt,
		Done:      t.Done,
		Id:        t.ID,
		Title:     t.Title,
	}
}

func toRespTaskList(ts []domain.Task) restapi.RespTaskList {
	l := make([]restapi.RespTask, 0, len(ts))
	for i := range ts {
		l = append(l, toRespTask(&ts[i]))
	}
	return restapi.RespTaskList{
		Items:      l,
		TotalCount: len(ts),
	}
}

func TestTaskAPI(t *testing.T) {
	tasks := []domain.Task{
		{
			UserID: gUser.ID,
			Title:  "Title0",
		},
		{
			UserID: gUser.ID,
			Title:  "Title1",
		},
		{
			UserID: gUser.ID,
			Title:  "Title2",
		},
		// 最後の要素だけ admin のタスク
		{
			UserID: gAdmin.ID,
			Title:  "Title4",
		},
	}

	t.Run("Create-OK", func(t *testing.T) {
		for i := range tasks {
			in := &tasks[i]
			t.Logf("Create %02d: UserID=%d", i, in.UserID)

			var token domain.AuthToken
			if in.UserID == gAdmin.ID {
				token = gAdmin.Token
			} else {
				token = gUser.Token
			}

			resp, err := gClient.CreateTask(gCtx, restapi.ReqCreateTask{
				Title: in.Title,
			}, addAuthTokenCookie(token))
			require.NoError(t, err)
			defer closeBody(resp)

			assert.Equal(t, 200, resp.StatusCode)
			got := decodeRespAsJSON[restapi.RespTask](t, resp)
			want := restapi.RespTask{
				CreatedAt: gClock.Now(),
				Done:      false,
				Id:        got.Id,
				Title:     in.Title,
			}
			assert.Equal(t, want, got)
			tasks[i].ID = got.Id
			tasks[i].CreatedAt = got.CreatedAt
		}
	})
	t.Run("List-OK", func(t *testing.T) {
		params := &restapi.ListTasksParams{Status: nil}
		resp, err := gClient.ListTasks(gCtx, params, addAuthTokenCookie(gUser.Token))
		require.NoError(t, err)
		defer closeBody(resp)

		assert.Equal(t, 200, resp.StatusCode)
		got := decodeRespAsJSON[restapi.RespTaskList](t, resp)
		want := toRespTaskList(tasks[:len(tasks)-1]) // tasksの最後の要素はadminのの所有物なので-1する
		assert.Equal(t, want, got)
	})
	t.Run("Patch-OK", func(t *testing.T) {
		targetTask := &tasks[0]
		req := restapi.ReqPatchTask{
			Done:  alloc(true),
			Title: nil,
		}
		targetTask.Done = *req.Done
		resp, err := gClient.PatchTask(gCtx, targetTask.ID, req, addAuthTokenCookie(gUser.Token))
		require.NoError(t, err)
		defer closeBody(resp)

		assert.Equal(t, 200, resp.StatusCode)
		got := decodeRespAsJSON[restapi.RespTask](t, resp)
		want := toRespTask(targetTask)
		assert.Equal(t, want, got)
	})
	t.Run("Get-OK", func(t *testing.T) {
		targetTask := &tasks[0]
		resp, err := gClient.GetTask(gCtx, targetTask.ID, addAuthTokenCookie(gUser.Token))
		require.NoError(t, err)
		defer closeBody(resp)

		assert.Equal(t, 200, resp.StatusCode)
		got := decodeRespAsJSON[restapi.RespTask](t, resp)
		want := toRespTask(targetTask)
		assert.Equal(t, want, got)
	})

	delTask := removeAt(1, &tasks)

	t.Run("Delete", func(t *testing.T) {
		resp, err := gClient.DeleteTask(gCtx, delTask.ID, addAuthTokenCookie(gUser.Token))
		require.NoError(t, err)
		defer closeBody(resp)

		assert.Equal(t, 200, resp.StatusCode)
		got := decodeRespAsJSON[restapi.RespTask](t, resp)
		want := toRespTask(&delTask)
		assert.Equal(t, want, got)
	})
	t.Run("Get-Deleted_task_NotFound", func(t *testing.T) {
		resp, err := gClient.GetTask(gCtx, delTask.ID, addAuthTokenCookie(gUser.Token))
		require.NoError(t, err)
		defer closeBody(resp)

		assert.Equal(t, 404, resp.StatusCode)
		got := decodeRespAsSimpleErrorJSON(t, resp)
		assert.Contains(t, got.Message, domain.ErrNotFound.Error())
	})
	t.Run("List-OnlyDone", func(t *testing.T) {
		params := &restapi.ListTasksParams{
			Status: alloc(restapi.TaskStatusFilterDone),
		}
		resp, err := gClient.ListTasks(gCtx, params, addAuthTokenCookie(gUser.Token))
		require.NoError(t, err)
		defer closeBody(resp)

		assert.Equal(t, 200, resp.StatusCode)
		got := decodeRespAsJSON[restapi.RespTaskList](t, resp)
		want := toRespTaskList(tasks[:1]) // tasks[0]しかdoneのものはない
		assert.Equal(t, want, got)
	})
}
