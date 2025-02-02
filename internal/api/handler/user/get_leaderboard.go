package user_api

import (
	"crypto_scam/internal/converter"
	"crypto_scam/internal/logger"
	"encoding/json"
	"net/http"
)

// GetLeaderboardHandler отправляет на клиент
// список 100 лучших пользователей в каждой лиге.
//
// В результате выполнения функции на клиент отправляется json с
// данными о пользователе, представленными структурой GetLeaderboardResponse.
func GetLeaderboardHandler(w http.ResponseWriter, r *http.Request) {
	// получаем данные из списка лидеров
	leaders, err := database.SelectLeaders()
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
