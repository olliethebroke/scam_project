package converter

import (
	modelAPI "crypto_scam/internal/api/handler/admin/model"
	modelDB "crypto_scam/internal/repository/model"
)

// TasksToGetTasksResponse конвертирует слайс заданий приходящий из бд,
// содержащий указатели на тип данных Task, в слайс, необходимый для отправки на клиент,
// содержащий указатели на тип данных GetTasksResponse.
//
// Входным значением для функции является слайс указателей на Task.
//
// Выходным значением для функции является слайс указателей на GetTasksResponse.
func TasksToGetTasksResponse(from []*modelDB.Task) []*modelAPI.GetTasksResponse {
	// инициализируем слайс указателей на структуру GetTasksResponse
	// такого же размера, как и слайс from
	to := make([]*modelAPI.GetTasksResponse, 0, len(from))

	// проходимся по каждому заданию из пришедшего слайса
	for _, task := range from {
		// и создаём указатель на структуру GetTasksResponse
		// с теми же значениями
		responseTask := &modelAPI.GetTasksResponse{
			Id:          task.Id,
			Reward:      task.Reward,
			Description: task.Description,
			ActionType:  task.ActionType,
			ActionData:  task.ActionData,
		}

		// добавляем указатель в выходной слайс
		to = append(to, responseTask)
	}

	// возвращаем слайс с заданиями
	return to
}
