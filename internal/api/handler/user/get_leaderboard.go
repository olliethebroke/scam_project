package user_api

import (
	"crypto_scam/internal/converter"
	"crypto_scam/internal/logger"
	"crypto_scam/internal/repository"
	"encoding/json"
	"net/http"
)

// GetLeaderboardHandler отправляет на клиент
// список 100 лучших пользователей в каждой лиге.
//
// В результате выполнения функции на клиент отправляется json с
// данными о пользователе, представленными структурой GetLeaderboardResponse.
func GetLeaderboardHandler(w http.ResponseWriter, r *http.Request) {
	// получаем контекст запроса
	ctx := r.Context()

	// получаем реализацию интерфейса Repository из контекста.
	// если значение отсутствует или имеет неверный тип, возвращаем ошибку 500.
	db, ok := ctx.Value("db").(repository.Repository)
	if !ok {
		http.Error(w, "failed to get database from context", http.StatusInternalServerError)
		logger.Warn("get_user.go/GetLeaderboardHandler - error while getting database from context")
		return
	}
	// получаем данные из списка лидеров
	leaders, err := db.SelectLeaders(ctx)
	if err != nil {
		http.Error(w, "failed to get leaders", http.StatusInternalServerError)
		logger.Warn("handler.go/GetLeaderboardHandler - error while selecting leaders: ", err)
		return
	}

	// энкодим данные о лидерах
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(converter.LeadersToGetLeaderboardResponse(leaders)); err != nil {
		http.Error(w, "failed to encode leaders", http.StatusInternalServerError)
		logger.Warn("handler.go/GetLeaderboardHandler - error while encoding data: ", err)
		return
	}
}
