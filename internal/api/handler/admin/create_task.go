package admin_api

import (
	"crypto_scam/internal/api/handler/admin/model"
	"crypto_scam/internal/converter"
	"crypto_scam/internal/logger"
	"encoding/json"
	"net/http"
)

// CreateTaskHandler создаёт новое задание, принимая информацию о нём в json формате.

// CreateTaskHandler получает от клиента информацию о
// новом задании и добавляет её в базу данных.
//
// В качестве входных данных выступает json, содержащий поля
// с описанием задания.
//
// В результате выполнения функции на клиент отправляется
// статус выполнения запроса.
func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	// декодим информацию о задании
	taskToCreate := &model.CreateTaskRequest{}
	err := json.NewDecoder(r.Body).Decode(taskToCreate)
	if err != nil {
		http.Error(w, "invalid json format", http.StatusBadRequest)
		logger.Warn("handler.go/CreateTaskHandler - error while decoding taskToCreate: ", err)
		return
	}

	// добавляем полученное задание в бд
	err = database.InsertTask(converter.CreateTaskRequestToTask(taskToCreate))
	if err != nil {
		http.Error(w, "failed to create task", http.StatusInternalServerError)
		logger.Warn("handler.go/CreateTaskHandler - error while inserting taskToCreate: ", err)
		return
	}

	// отправляем клиенту статус успешного создания задания
	w.WriteHeader(http.StatusCreated)
}
