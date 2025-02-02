package admin_api

import (
	"crypto_scam/internal/converter"
	"crypto_scam/internal/logger"
	"encoding/json"
	"net/http"
)

// GetTasksHandler отправляет на клиент список заданий из базы данных.
//
// В результате выполнения функции на клиент отправляется json с
// данными о всех игровых заданиях, представленными
// структурой GetTasksResponse.
func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	// получаем из базы данных задания
	tasks, err := database.SelectTasks()
	if err != nil {
		http.Error(w, "failed to get tasks", http.StatusNotFound)
		logger.Warn("get_tasks.go/GetTasksHandler - error while selecting tasks: ", err)
		return
	}
	// отправляем на клиент задания
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(converter.TasksToGetTasksResponse(tasks))
	if err != nil {
		http.Error(w, "failed to encode tasks", http.StatusInternalServerError)
		logger.Warn("get_tasks.go/GetTasksHandler - error while encoding tasks: ", err)
		return
	}
}
