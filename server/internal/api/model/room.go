package model

type UserRoom struct {
	Id       string `json:"id"`
	RoomId   string `json:"roomId" validate:"required"`
	UserId   string `json:"userId" validate:"required"`
	Username string `json:"username" validate:"required"`
}
