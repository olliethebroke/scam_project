package model

// UpdateUserRequest описывает данные для обновления
// игровых показателей пользователя в базе данных.
//
// Blocks - количество блоков,
// Record - игровой рекорд пользователя.
type UpdateUserRequest struct {
	Blocks int64 `json:"blocks"`
	Record int64 `json:"record"`
}
