// Package restapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.4 DO NOT EDIT.
package restapi

import (
	"time"

	"github.com/arumakan1727/todo-app-go-react/domain"
	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
)

const (
	AuthCookieScopes      = "AuthCookie.Scopes"
	BearerAuthTokenScopes = "BearerAuthToken.Scopes"
)

// Defines values for TaskStatusFilter.
const (
	TaskStatusFilterAny  TaskStatusFilter = "any"
	TaskStatusFilterDone TaskStatusFilter = "done"
	TaskStatusFilterTodo TaskStatusFilter = "todo"
)

// Password defines model for Password.
type Password = string

// ReqCreateAuthToken defines model for ReqCreateAuthToken.
type ReqCreateAuthToken struct {
	Email    openapi_types.Email `json:"email"`
	Password Password            `json:"password"`
}

// ReqCreateTask defines model for ReqCreateTask.
type ReqCreateTask struct {
	Title TaskTitle `json:"title"`
}

// ReqCreateUser defines model for ReqCreateUser.
type ReqCreateUser struct {
	DisplayName string              `json:"displayName"`
	Email       openapi_types.Email `json:"email"`
	Password    Password            `json:"password"`
}

// ReqPatchTask defines model for ReqPatchTask.
type ReqPatchTask struct {
	Done  *bool      `json:"done,omitempty"`
	Title *TaskTitle `json:"title,omitempty"`
}

// RespTask defines model for RespTask.
type RespTask struct {
	CreatedAt time.Time `json:"createdAt"`
	Done      bool      `json:"done"`
	Id        TaskID    `json:"id"`
	Title     string    `json:"title"`
}

// RespTaskList defines model for RespTaskList.
type RespTaskList struct {
	Items      []RespTask `json:"items"`
	TotalCount int        `json:"totalCount"`
}

// RespUser defines model for RespUser.
type RespUser struct {
	CreatedAt   time.Time           `json:"createdAt"`
	DisplayName string              `json:"displayName"`
	Email       openapi_types.Email `json:"email"`
	Id          UserID              `json:"id"`
}

// RespUserList defines model for RespUserList.
type RespUserList struct {
	Items      []RespUser `json:"items"`
	TotalCount int        `json:"totalCount"`
}

// TaskID defines model for TaskID.
type TaskID = domain.TaskID

// TaskStatusFilter タスクの完了状態のフィルタリング指定
type TaskStatusFilter string

// TaskTitle defines model for TaskTitle.
type TaskTitle = string

// UserID defines model for UserID.
type UserID = domain.UserID

// Resp200Task defines model for Resp200Task.
type Resp200Task = RespTask

// ListTasksParams defines parameters for ListTasks.
type ListTasksParams struct {
	// Status タスクの完了状態のフィルタリング指定
	Status *TaskStatusFilter `form:"status,omitempty" json:"status,omitempty"`
}

// CreateAuthTokenJSONRequestBody defines body for CreateAuthToken for application/json ContentType.
type CreateAuthTokenJSONRequestBody = ReqCreateAuthToken

// CreateTaskJSONRequestBody defines body for CreateTask for application/json ContentType.
type CreateTaskJSONRequestBody = ReqCreateTask

// PatchTaskJSONRequestBody defines body for PatchTask for application/json ContentType.
type PatchTaskJSONRequestBody = ReqPatchTask

// CreateUserJSONRequestBody defines body for CreateUser for application/json ContentType.
type CreateUserJSONRequestBody = ReqCreateUser
