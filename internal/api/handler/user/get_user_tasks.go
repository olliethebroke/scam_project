package user_api

import (
	"crypto_scam/internal/converter"
	"crypto_scam/internal/logger"
	"crypto_scam/internal/repository"
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
	// получаем контекст запроса
	ctx := r.Context()

	// получаем реализацию интерфейса Repository из контекста.
	// если значение отсутствует или имеет неверный тип, возвращаем ошибку 500.
	db, ok := ctx.Value("db").(repository.Repository)
	if !ok {
		http.Error(w, "failed to get database from context", http.StatusInternalServerError)
		logger.Warn("get_user.go/GetUserTasksHandler - error while getting database from context")
		return
	}
	// получаем id пользователя из контекста.
	// если значение отсутствует или имеет неверный тип, возвращаем ошибку 500.
	id, ok := ctx.Value("id").(int64)
	if !ok {
		http.Error(w, "failed to get user id", http.StatusInternalServerError)
		logger.Warn("get_user.go/GetUserTasksHandler - error while getting user id from context")
		return
	}

	// получаем задания из бд
	tasks, err := db.SelectUserTasks(ctx, id)
	if err != nil {
		http.Error(w, "failed to get user tasks", http.StatusInternalServerError)
		logger.Warn("handler.go/GetUserTasksHandler - error while selecting tasks: ", err)
		return
	}

	// энкодим данные о заданиях
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(converter.TasksToGetUserTasksResponse(tasks)); err != nil {
		http.Error(w, "failed to encode user tasks", http.StatusInternalServerError)
		logger.Warn("handler.go/GetUserTasksHandler - error while selecting tasks: ", err)
		return
	}
}
