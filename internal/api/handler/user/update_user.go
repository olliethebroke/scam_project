package user_api

import (
	"crypto_scam/internal/api/handler/user/model"
	"crypto_scam/internal/converter"
	"crypto_scam/internal/logger"
	"crypto_scam/internal/repository"
	"encoding/json"
	"net/http"
)

// UpdateUserHandler получает от клиента данные об
// игровых показателях пользователя
// и обновляет информацию в базе данных.
//
// В качестве входных данных выступает id пользователя из телеграма,
// который передаётся в контексте запроса и json, содержащий
// данные о количестве блоков и рекорде пользователя.
//
// В результате выполнения функции на клиент отправляется
// статус выполнения запроса.
func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	// получаем контекст запроса
	ctx := r.Context()

	// получаем реализацию интерфейса Repository из контекста.
	// если значение отсутствует или имеет неверный тип, возвращаем ошибку 500.
	db, ok := ctx.Value("db").(repository.Repository)
	if !ok {
		http.Error(w, "failed to get database from context", http.StatusInternalServerError)
		logger.Warn("get_user.go/UpdateUserHandler - error while getting database from context")
		return
	}
	// получаем id пользователя из контекста.
	// если значение отсутствует или имеет неверный тип, возвращаем ошибку 500.
	id, ok := ctx.Value("id").(int64)
	if !ok {
		http.Error(w, "failed to get user id", http.StatusInternalServerError)
		logger.Warn("get_user.go/UpdateUserHandler - error while getting user id from context")
		return
	}

	// декодим тело запроса, содержащее обновлённые данные о пользователе
	dataToUpdate := &model.UpdateUserRequest{}
	if err := json.NewDecoder(r.Body).Decode(dataToUpdate); err != nil {
		http.Error(w, "invalid json format", http.StatusBadRequest)
		logger.Warn("handler.go/UpdateUserHandler - error while decoding user data: ", err)
		return
	}

	// обновляем данные в бд
	err := db.UpdateUser(ctx, id, converter.UpdateUserRequestToUpdateUser(dataToUpdate))
	if err != nil {
		http.Error(w, "failed to update user", http.StatusInternalServerError)
		logger.Warn("handler.go/UpdateUserHandler - error while updating user: ", err)
		return
	}

	// уведомляем клиент об успешном завершении процесса
	w.WriteHeader(http.StatusOK)
}
