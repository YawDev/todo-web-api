package helpers

import "time"

type User struct {
	Username string `binding:"required"`
	Password string `binding:"required"`
}

type SuccessResponse struct {
	Status  int    `json:"status" example:"200"`
	Message string `json:"message" example:"Success"`
}

type ErrorResponse struct {
	Status  int    `json:"status" example:"500"`
	Message string `json:"message" example:"Something went wrong"`
}

type BadRequestResponse struct {
	Status  int    `json:"status" example:"400"`
	Message string `json:"message" example:"Bad request"`
}

type NotFoundResponse struct {
	Status  int    `json:"status" example:"404"`
	Message string `json:"message" example:"Not Found"`
}

type UnauthorizedResponse struct {
	Status  int    `json:"status" example:"401"`
	Message string `json:"message" example:"Unauthorized"`
}

type SaveResponse struct {
	Status  int    `json:"status" example:"200"`
	Message string `json:"message" example:"Successfully saved"`
	Id      int    `json:"id" example:"1"`
}

type DeleteResult struct {
	Status  int    `json:"status" example:"200"`
	Message string `json:"message" example:"Successfully deleted"`
	Success bool   `json:"success" example:"true"`
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
