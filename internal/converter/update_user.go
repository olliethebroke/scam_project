package converter

import (
	modelAPI "crypto_scam/internal/api/handler/user/model"
	modelDB "crypto_scam/internal/repository/model"
)

// UpdateUserRequestToUpdateUser конвертирует данные
// из структуры клиентского запроса UpdateUserRequest
// в структуру Update.
//
// Входным параметром для функции является указатель на структуру UpdateUserRequest.
//
// Выходным параметром для функции является указатель на структуру Update.
func UpdateUserRequestToUpdateUser(from *modelAPI.UpdateUserRequest) *modelDB.Update {
	// создаём и инициализируем указатель на структуру Update
	// заполняем нужными значениями
	to := &modelDB.Update{
		Blocks: from.Blocks,
		Record: from.Record,
	}

	// возвращаем указатель на Update
	return to
}
