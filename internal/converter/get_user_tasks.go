package converter

import (
	modelAPI "crypto_scam/internal/api/handler/user/model"
	modelDB "crypto_scam/internal/repository/model"
)

// TasksToGetUserTasksResponse конвертирует задания пользователя из слайса указателей
// на тип данных Task в слайс указателей на тип данных GetUserTasksResponse.
//
// Входным значением для функции является слайс указателей на Task.
//
// Выходным значением для функции является слайс указателей на GetUserTasksResponse.
func TasksToGetUserTasksResponse(from []*modelDB.Task) []*modelAPI.GetUserTasksResponse {
	// инициализируем слайс указателей на структуру GetTasksResponse
	// такого же размера, как и слайс from
	to := make([]*modelAPI.GetUserTasksResponse, 0, len(from))

	// проходимся по каждому заданию из пришедшего слайса
	for _, task := range from {
		// и создаём указатель на структуру GetTasksResponse
		// с теми же значениями
		responseTask := &modelAPI.GetUserTasksResponse{
			Id:          task.Id,
			Reward:      task.Reward,
			Description: task.Description,
			ActionType:  task.ActionType,
			ActionData:  task.ActionData,
			IsCompleted: task.IsCompleted,
		}

		// добавляем указатель в выходной слайс
		to = append(to, responseTask)
	}

	// возвращаем слайс с заданиями
	return to
}
