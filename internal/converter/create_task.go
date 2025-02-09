package converter

import (
	modelAPI "crypto_scam/internal/api/handler/admin/model"
	modelDB "crypto_scam/internal/repository/model"
)

// CreateTaskRequestToTask преобразует данные из структуры CreateTaskRequest
// в структуру Task.
// Используется для конвертации данных с клиента в данные для работы с базой данных.
//
// Входные данные для функции - указатель на тип данных CreateTaskRequest.
//
// Выходные данные - указатель на тип данных Task.
func CreateTaskRequestToTask(from *modelAPI.CreateTaskRequest) *modelDB.Task {
	// создаём переменную, которая указывает на структуру Task
	to := &modelDB.Task{}

	// заполняем структу данными
	to.Reward = from.Reward
	to.ActionData = from.ActionData
	to.ActionType = from.ActionType
	to.Description = from.Description

	// возвращаем указатель
	return to
}

// TaskToCreateTaskResponse преобразует данные из структуры Task
// в структуру CreateTaskResponse.
// Используется для конвертации данных из бд в данные для ответа клиенту.
//
// Входные данные для функции - указатель на тип данных Task.
//
// Выходные данные - указатель на тип данных CreateTaskResponse.
func TaskToCreateTaskResponse(from *modelDB.Task) *modelAPI.CreateTaskResponse {
	to := &modelAPI.CreateTaskResponse{
		Id: from.Id,
	}
	return to
}
