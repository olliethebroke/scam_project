package user_api

import (
	"crypto_scam/internal/api/handler/user/model"
	"crypto_scam/internal/converter"
	"crypto_scam/internal/logger"
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
	// получаем id пользователя из контекста запроса
	id, ok := r.Context().Value("id").(int64)
	if !ok {
		http.Error(w, "failed to get user id", http.StatusInternalServerError)
		logger.Warn("update_user.go/UpdateUserHandler - error while getting user id from context")
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
	err := database.UpdateUser(id, converter.UpdateUserRequestToUpdateUser(dataToUpdate))
	if err != nil {
		http.Error(w, "failed to update user", http.StatusInternalServerError)
		logger.Warn("handler.go/UpdateUserHandler - error while updating user: ", err)
		return
	}

	// уведомляем клиент об успешном завершении процесса
	w.WriteHeader(http.StatusOK)
}
