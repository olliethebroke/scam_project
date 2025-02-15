package admin_api

import (
	"crypto_scam/internal/api/handler/admin/model"
	"crypto_scam/internal/converter"
	"crypto_scam/internal/logger"
	"crypto_scam/internal/repository"
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
	// получаем контекст запроса
	ctx := r.Context()

	// получаем реализацию интерфейса Repository из контекста.
	// если значение отсутствует или имеет неверный тип, возвращаем ошибку 500.
	db, ok := ctx.Value("db").(repository.Repository)
	if !ok {
		http.Error(w, "failed to get database from context", http.StatusInternalServerError)
		logger.Warn("get_user.go/CreateTaskHandler - error while getting database from context")
		return
	}
	// декодим информацию о задании
	taskToCreate := &model.CreateTaskRequest{}
	err := json.NewDecoder(r.Body).Decode(taskToCreate)
	if err != nil {
		http.Error(w, "invalid json format", http.StatusBadRequest)
		logger.Warn("handler.go/CreateTaskHandler - error while decoding task: ", err)
		return
	}

	// добавляем полученное задание в бд
	createdTask, err := db.InsertTask(ctx, converter.CreateTaskRequestToTask(taskToCreate))
	if err != nil {
		http.Error(w, "failed to create task", http.StatusInternalServerError)
		logger.Warn("handler.go/CreateTaskHandler - error while inserting task: ", err)
		return
	}

	// кодируем данные и отправляем на клиент
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(converter.TaskToCreateTaskResponse(createdTask))
	if err != nil {
		http.Error(w, "failed to encode task", http.StatusInternalServerError)
		logger.Warn("handler.go/CreateTaskHandler - error while encoding task: ", err)
		return
	}
}
