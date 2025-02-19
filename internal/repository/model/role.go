package model

// UserRole структура, описывающая роль игрока,
// которому предоставляется доступ к апи.
//
// Id - идентификатор пользователя в телеграм,
// Role - роль пользователя.
type UserRole struct {
	Id   int64
	Role int16
}
