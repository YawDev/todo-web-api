package helpers

import "time"

type User struct {
	Username string `binding:"required"`
	Password string `binding:"required"`
}

type SuccessResponse struct {
	Status  string `json:"status" example:"200"`
	Message string `json:"message" example:"success"`
}

type ErrorResponse struct {
	Status  int    `json:"status" example:"500"`
	Message string `json:"message" example:"Something went wrong"`
}

type BadRequestResponse struct {
	Status  int    `json:"status" example:"400"`
	Message string `json:"message" example:"Bad request"`
}

type SaveResponse struct {
	Status  int    `json:"status" example:"200"`
	Message string `json:"message" example:"Successfully created`
	Id      int    `json:"id" example:"1`
}

type DeleteResult struct {
	Status  int    `json:"status" example:"200"`
	Message string `json:"message" example:"Successfully deleted"`
	Success bool   `json:"success" example:"true`
}

type UserResult struct {
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"createdAt"`
}

type SaveTask struct {
	Title       string `binding:"required"`
	Description string
}

type SetStatus struct {
	IsCompleted bool `binding:"required"`
}
