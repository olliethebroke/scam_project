package model

// DeleteUserRequest описывает пользователя,
// которого нужно удалить.
//
// Id - идентификатор пользователя.
type DeleteUserRequest struct {
	Id int64 `json:"id"`
}
