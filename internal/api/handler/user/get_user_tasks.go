package user_api

import (
	"crypto_scam/internal/logger"
	"encoding/json"
	"net/http"
)

// GetUserTasksHandler отправляет на клиент информацию о заданиях
// пользователя, сделавшего запрос.
//
// В качестве входных данных выступает id пользователя из телеграма,
// который передаётся в контексте запроса.
//
// В результате выполнения функции на клиент отправляется json с
// данными о заданиях конкретного пользователя,
// представленными структурой GetUserTasksResponse.
func GetUserTasksHandler(w http.ResponseWriter, r *http.Request) {
	// получаем id пользователя из контекста запроса
	id, ok := r.Context().Value("id").(int64)
	if !ok {
		http.Error(w, "failed to get user id", http.StatusInternalServerError)
		logger.Warn("get_user_tasks.go/GetUserTasksHandler - error while getting user id from context")
		return
	}

	// получаем задания из бд
	tasks, err := database.SelectUserTasks(id)
	if err != nil {
		http.Error(w, "failed to get user tasks", http.StatusInternalServerError)
		logger.Warn("handler.go/GetUserTasksHandler - error while selecting tasks: ", err)
		return
	}

	// энкодим данные о заданиях
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(tasks); err != nil {
		http.Error(w, "failed to encode user tasks", http.StatusInternalServerError)
		logger.Warn("handler.go/GetUserTasksHandler - error while selecting tasks: ", err)
		return
	}
}
