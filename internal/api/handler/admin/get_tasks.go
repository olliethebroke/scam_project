package admin_api

import (
	"crypto_scam/internal/converter"
	"crypto_scam/internal/logger"
	"crypto_scam/internal/repository"
	"encoding/json"
	"net/http"
)

// GetTasksHandler отправляет на клиент список заданий из базы данных.
//
// В результате выполнения функции на клиент отправляется json с
// данными о всех игровых заданиях, представленными
// структурой GetTasksResponse.
func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	// получаем реализацию интерфейса Repository из контекста.
	// если значение отсутствует или имеет неверный тип, возвращаем ошибку 500.
	db, ok := r.Context().Value("db").(repository.Repository)
	if !ok {
		http.Error(w, "failed to get database from context", http.StatusInternalServerError)
		logger.Warn("get_user.go/GetTasksHandler - error while getting database from context")
		return
	}
	// получаем из базы данных задания
	tasks, err := db.SelectTasks()
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
