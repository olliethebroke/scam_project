package model

// CreateUserRequest содержит поле запроса от клиента на создание пользователя
type CreateUserRequest struct {
	Username string `json:"username"`
}
