package api

import (
	"crypto_scam/internal/logger/logrus_logger"
	"crypto_scam/internal/storage/postgres"
	"crypto_scam/pkg/utils/string_utils"
	"encoding/json"
	"github.com/go-chi/chi"
	"net/http"
)

// CreateUserHandler парсит id из телеграма и отправляет информацию о созданном пользователе
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	logrus_logger.Log.Infof("received request %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
	userId := chi.URLParam(r, "id")
	id, err := string_utils.ParseID(userId)
	if err != nil {
		http.Error(w, "internal/api/handler.go - error while parsing user id", http.StatusBadRequest)
		logrus_logger.Log.Warn("internal/api/handler.go - error while parsing user id: ", err)
		return
	}
	var req postgres.CreateRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "internal/api/handler.go - error while decoding username", http.StatusBadRequest)
		logrus_logger.Log.Warn("internal/api/handler.go - error while decoding username: ", err)
		return
	}
	info, err := postgres.InsertUser(id, req.Username)
	if err != nil {
		http.Error(w, "internal/api/handler.go - error while inserting user", http.StatusInternalServerError)
		logrus_logger.Log.Warn("internal/api/handler.go - error while inserting user: ", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(info); err != nil {
		http.Error(w, "internal/api/handler.go - error while encoding data", http.StatusInternalServerError)
		logrus_logger.Log.Warn("internal/api/handler.go - error while encoding data: ", err)
		return
	}
}

// GetUserHandler парсит id из телеграма и отправляет информацию о пользователе
func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	logrus_logger.Log.Infof("received request %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
	userId := chi.URLParam(r, "id")
	id, err := string_utils.ParseID(userId)
	if err != nil {
		http.Error(w, "internal/api/handler.go - error while parsing user id", http.StatusBadRequest)
		logrus_logger.Log.Warn("internal/api/handler.go - error while parsing user id: ", err)
		return
	}
	info, err := postgres.SelectUser(id)
	if err != nil {
		http.Error(w, "internal/api/handler.go - error while selecting user", http.StatusInternalServerError)
		logrus_logger.Log.Warn("internal/api/handler.go - error while selecting user: ", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(info); err != nil {
		http.Error(w, "internal/api/handler.go - error while encoding data", http.StatusInternalServerError)
		logrus_logger.Log.Warn("internal/api/handler.go - error while encoding data: ", err)
		return
	}
}

// GetLeaderboardHandler возвращает список 100 самых лучших пользователей в каждой лиге
func GetLeaderboardHandler(w http.ResponseWriter, r *http.Request) {
	logrus_logger.Log.Infof("received request %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
	leaders, err := postgres.SelectLeaders()
	if err != nil {
		http.Error(w, "internal/api/handler.go - error while selecting leaders", http.StatusInternalServerError)
		logrus_logger.Log.Warn("internal/api/handler.go - error while selecting leaders: ", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(leaders); err != nil {
		http.Error(w, "internal/api/handler.go - error while encoding data", http.StatusInternalServerError)
		logrus_logger.Log.Warn("internal/api/handler.go - error while encoding data: ", err)
		return
	}
}

// UpdateUserHandler обновляет значения пользователя с помощью полученных данных
func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	logrus_logger.Log.Infof("received request %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
	userId := chi.URLParam(r, "id")
	id, err := string_utils.ParseID(userId)
	if err != nil {
		http.Error(w, "internal/api/handler.go - error while parsing user id", http.StatusBadRequest)
		logrus_logger.Log.Warn("internal/api/handler.go - error while parsing user id: ", err)
		return
	}
	info := &postgres.UpdateRequest{}
	if err = json.NewDecoder(r.Body).Decode(info); err != nil {
		http.Error(w, "internal/api/handler.go - error while decoding user data", http.StatusBadRequest)
		logrus_logger.Log.Warn("internal/api/handler.go - error while decoding user data: ", err)
		return
	}
	err = postgres.UpdateUser(id, info)
	if err != nil {
		http.Error(w, "internal/api/handler.go - error while updating user", http.StatusInternalServerError)
		logrus_logger.Log.Warn("internal/api/handler.go - error while updating user: ", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
