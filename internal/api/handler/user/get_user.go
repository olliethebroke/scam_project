package user_api

import (
	"crypto_scam/internal/converter"
	"crypto_scam/internal/logger"
	"encoding/json"
	"net/http"
)

// GetUserHandler отправляет на клиент информацию о пользователе,
// сделавшем запрос.
//
// В качестве входных данных выступает id пользователя из телеграма, который
// передаётся в контексте запроса.
//
// В результате выполнения функции на клиент отправляется json с
// данными о пользователе, представленными структурой GetUserResponse.
func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	// получаем id пользователя из контекста запроса
	id, ok := r.Context().Value("id").(int64)
	if !ok {
		http.Error(w, "failed to get user id", http.StatusInternalServerError)
		logger.Warn("get_user.go/GetUserHandler - error while getting user id from context")
		return
	}

	// получаем данные о пользователе из бд
	user, err := database.SelectUser(id)
	if err != nil {
		http.Error(w, "failed to get user", http.StatusNotFound)
		logger.Warn("handler.go/GetUserHandler - error while selecting user: ", err)
		return
	}

	// энкодим данные о пользователе
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(converter.UserToGetUserResponse(user)); err != nil {
		http.Error(w, "failed to encode user", http.StatusInternalServerError)
		logger.Warn("handler.go/GetUserHandler - error while encoding data: ", err)
		return
	}
}
