package model

// CreateUserRequest содержит поле запроса от клиента на создание пользователя
type CreateUserRequest struct {
	Id        int64  `json:"id"`
	Firstname string `json:"firstname"`
}
