package admin_api

import (
	"crypto_scam/internal/api/handler/admin/model"
	"crypto_scam/internal/logger"
	"crypto_scam/internal/repository"
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
	// получаем контекст запроса
	ctx := r.Context()

	// получаем реализацию интерфейса Repository из контекста.
	// если значение отсутствует или имеет неверный тип, возвращаем ошибку 500.
	db, ok := ctx.Value("db").(repository.Repository)
	if !ok {
		http.Error(w, "failed to get database from context", http.StatusInternalServerError)
		logger.Warn("get_user.go/DeleteUserHandler - error while getting database from context")
		return
	}
	userToDelete := &model.DeleteUserRequest{}
	err := json.NewDecoder(r.Body).Decode(userToDelete)
	if err != nil {
		http.Error(w, "invalid json format", http.StatusBadRequest)
		logger.Warn("delete_user.go/DeleteUserHandler - error while decoding user id: ", err)
		return
	}

	// удаляем пользователя из бд
	err = db.DeleteUser(ctx, userToDelete.Id)
	if err != nil {
		http.Error(w, "failed to delete user", http.StatusInternalServerError)
		logger.Warn("delete_user.go/DeleteUserHandler - error while deleting user: ", err)
		return
	}

	// если не возникло ошибок
	// отправляем клиенту статус ОК
	w.WriteHeader(http.StatusOK)
}
