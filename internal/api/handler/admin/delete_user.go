package admin_api

import (
	"crypto_scam/internal/api/handler/admin/model"
	"crypto_scam/internal/logger"
	"encoding/json"
	"net/http"
)

// DeleteUserHandler удаляет пользователя из базы данных.
//
// В качестве входных данных функция принимает json,
// содержащий поле с идентификатором пользователя, которого
// необходимо удалить.
//
// В результате выполнения функции на клиент отправляется
// статус выполнения запроса.
func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	userToDelete := &model.DeleteUserRequest{}
	err := json.NewDecoder(r.Body).Decode(userToDelete)
	if err != nil {
		http.Error(w, "invalid json format", http.StatusBadRequest)
		logger.Warn("delete_user.go/DeleteUserHandler - error while decoding user id: ", err)
		return
	}

	// удаляем пользователя из бд
	err = database.DeleteUser(userToDelete.Id)
	if err != nil {
		http.Error(w, "failed to delete user", http.StatusInternalServerError)
		logger.Warn("delete_user.go/DeleteUserHandler - error while deleting user: ", err)
		return
	}

	// если не возникло ошибок
	// отправляем клиенту статус ОК
	w.WriteHeader(http.StatusOK)
}
