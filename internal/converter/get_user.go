package converter

import (
	modelAPI "crypto_scam/internal/api/handler/user/model"
	modelDB "crypto_scam/internal/repository/model"
)

// UserToGetUserResponse конвертирует данные о пользователе из
// типа данных User в тип данных GetUserResponse
// для отправки на клиент.
//
// Входным параметром функции является указатель на структуру User.
//
// Выходным параметром функции является указатель на структуру GetUserResponse
func UserToGetUserResponse(from *modelDB.User) *modelAPI.GetUserResponse {
	// создаём и инициализируем указатель на структуру GetUserResponse
	to := &modelAPI.GetUserResponse{
		Id:             from.Id,
		Firstname:      from.Firstname,
		Blocks:         from.Blocks,
		Record:         from.Record,
		DaysStreak:     from.DaysStreak,
		InvitedFriends: from.InvitedFriends,
		IsPremium:      from.IsPremium,
		Award:          from.Award,
		League:         from.League,
	}

	// возвращаем указатель на GetUserResponse
	return to
}
