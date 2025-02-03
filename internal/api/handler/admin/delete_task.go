package admin_api

import (
	"crypto_scam/internal/api/handler/admin/model"
	"crypto_scam/internal/logger"
	"crypto_scam/internal/repository"
	"encoding/json"
	"net/http"
)

// DeleteTaskHandler удаляет задание из базы данных.
//
// В качестве входных данных функция принимает json,
// содержащий поле с идентификатором задания, которое
// необходимо удалить.
//
// В результате выполнения функции на клиент отправляется
// статус выполнения запроса.
func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	// получаем реализацию интерфейса Repository из контекста.
	// если значение отсутствует или имеет неверный тип, возвращаем ошибку 500.
	db, ok := r.Context().Value("db").(repository.Repository)
	if !ok {
		http.Error(w, "failed to get database from context", http.StatusInternalServerError)
		logger.Warn("get_user.go/DeleteTaskHandler - error while getting database from context")
		return
	}
	// декодим полученный запрос
	taskToDelete := &model.DeleteTaskRequest{}
	err := json.NewDecoder(r.Body).Decode(taskToDelete)
	if err != nil {
		http.Error(w, "invalid json format", http.StatusBadRequest)
		logger.Warn("delete_task.go/DeleteTaskHandler - error while decoding task id: ", err)
		return
	}

	// удаляем задание из базы данных
	err = db.DeleteTask(taskToDelete.Id)
	if err != nil {
		http.Error(w, "failed to delete task", http.StatusInternalServerError)
	}

	// если всё выполнилось без ошибок
	// отправляем статус ОК
	w.WriteHeader(http.StatusOK)
}
