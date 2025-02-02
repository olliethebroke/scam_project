package admin_api

import (
	"crypto_scam/internal/api/handler/admin/model"
	"crypto_scam/internal/logger"
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
	// декодим полученный запрос
	taskToDelete := &model.DeleteTaskRequest{}
	err := json.NewDecoder(r.Body).Decode(taskToDelete)
	if err != nil {
		http.Error(w, "invalid json format", http.StatusBadRequest)
		logger.Warn("delete_task.go/DeleteTaskHandler - error while decoding task id: ", err)
		return
	}

	// удаляем задание из базы данных
	err = database.DeleteTask(taskToDelete.Id)
	if err != nil {
		http.Error(w, "failed to delete task", http.StatusInternalServerError)
	}

	// если всё выполнилось без ошибок
	// отправляем статус ОК
	w.WriteHeader(http.StatusOK)
}
