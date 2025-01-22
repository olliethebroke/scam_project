package api

import (
	"crypto_scam/internal/logger/logrus_logger"
	"crypto_scam/internal/storage/models"
	"crypto_scam/internal/storage/postgres"
	"crypto_scam/pkg/utils/string_utils"
	"encoding/json"
	"github.com/go-chi/chi"
	"net/http"
)

// CreateUserHandler парсит id и отправляет информацию о созданном пользователе
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	logrus_logger.Log.Infof("received request %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
	// получаем id пользователя из запроса
	userId := chi.URLParam(r, "id")
	// парсим его в int64
	id, err := string_utils.ParseID(userId)
	if err != nil {
		http.Error(w, "handlers.go/CreateUserHandler - error while parsing user id", http.StatusBadRequest)
		logrus_logger.Log.Warn("handlers.go/CreateUserHandler - error while parsing user id: ", err)
		return
	}
	var req models.CreateUserRequest
	// декодим тело запроса в структуру CreateUserRequest
	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "handlers.go/CreateUserHandler - error while decoding username", http.StatusBadRequest)
		logrus_logger.Log.Warn("handlers.go/CreateUserHandler - error while decoding username: ", err)
		return
	}
	// добавляем пользователя в бд и принимаем информацию о нём
	info, err := postgres.InsertUser(id, req.Username)
	if err != nil {
		http.Error(w, "handlers.go/CreateUserHandler - error while inserting user", http.StatusInternalServerError)
		logrus_logger.Log.Warn("handlers.go/CreateUserHandler - error while inserting user: ", err)
		return
	}
	// энкодим данные о созданном пользователе
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(info); err != nil {
		http.Error(w, "handlers.go/CreateUserHandler - error while encoding data", http.StatusInternalServerError)
		logrus_logger.Log.Warn("handlers.go/CreateUserHandler - error while encoding data: ", err)
		return
	}
}

// GetUserHandler парсит id из телеграма и отправляет информацию о пользователе
func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	logrus_logger.Log.Infof("received request %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
	// получаем id пользователя из запроса
	userId := chi.URLParam(r, "id")
	// парсим его в int64
	id, err := string_utils.ParseID(userId)
	if err != nil {
		http.Error(w, "handlers.go/GetUserHandler - error while parsing user id", http.StatusBadRequest)
		logrus_logger.Log.Warn("handlers.go/GetUserHandler - error while parsing user id: ", err)
		return
	}
	// получаем данные о пользователе из бд
	info, err := postgres.SelectUser(id)
	if err != nil {
		http.Error(w, "handlers.go/GetUserHandler - error while selecting user", http.StatusNotFound)
		logrus_logger.Log.Warn("handlers.go/GetUserHandler - error while selecting user: ", err)
		return
	}
	// энкодим данные о пользователе
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(info); err != nil {
		http.Error(w, "handlers.go/GetUserHandler - error while encoding data", http.StatusInternalServerError)
		logrus_logger.Log.Warn("handlers.go/GetUserHandler - error while encoding data: ", err)
		return
	}
}

// GetLeaderboardHandler возвращает список 100 самых лучших пользователей в каждой лиге
func GetLeaderboardHandler(w http.ResponseWriter, r *http.Request) {
	logrus_logger.Log.Infof("received request %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
	// получаем данные из списка лидеров
	leaders, err := postgres.SelectLeaders()
	if err != nil {
		http.Error(w, "handlers.go/GetLeaderboardHandler - error while selecting leaders", http.StatusInternalServerError)
		logrus_logger.Log.Warn("handlers.go/GetLeaderboardHandler - error while selecting leaders: ", err)
		return
	}
	// энкодим данные о лидерах
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(leaders); err != nil {
		http.Error(w, "handlers.go/GetLeaderboardHandler - error while encoding data", http.StatusInternalServerError)
		logrus_logger.Log.Warn("handlers.go/GetLeaderboardHandler - error while encoding data: ", err)
		return
	}
}

// UpdateUserHandler обновляет значения пользователя с помощью полученных данных
func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	logrus_logger.Log.Infof("received request %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
	// получаем id пользователя из запроса
	userId := chi.URLParam(r, "id")
	// парсим его в int64
	id, err := string_utils.ParseID(userId)
	if err != nil {
		http.Error(w, "handlers.go/UpdateUserHandler - error while parsing user id", http.StatusBadRequest)
		logrus_logger.Log.Warn("handlers.go/UpdateUserHandler - error while parsing user id: ", err)
		return
	}
	info := &models.UpdateUserRequest{}
	// декодим тело запроса, содержащее обновлённые данные о пользователе
	if err = json.NewDecoder(r.Body).Decode(info); err != nil {
		http.Error(w, "handlers.go/UpdateUserHandler - error while decoding user data", http.StatusBadRequest)
		logrus_logger.Log.Warn("handlers.go/UpdateUserHandler - error while decoding user data: ", err)
		return
	}
	// обновляем данные в бд
	err = postgres.UpdateUser(id, info)
	if err != nil {
		http.Error(w, "handlers.go/UpdateUserHandler - error while updating user", http.StatusInternalServerError)
		logrus_logger.Log.Warn("handlers.go/UpdateUserHandler - error while updating user: ", err)
		return
	}
	// уведомляем клиент об успешном завершении процесса
	w.WriteHeader(http.StatusOK)
}

// CreateFriendshipHandler создаёт запись о дружбе пользователей
func CreateFriendshipHandler(w http.ResponseWriter, r *http.Request) {
	logrus_logger.Log.Infof("received request %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
	friendship := &models.Friendship{}
	// декодим json с двумя id пользователей
	if err := json.NewDecoder(r.Body).Decode(friendship); err != nil {
		http.Error(w, "handlers.go/CreateFriendshipHandler - error while decoding friendship", http.StatusBadRequest)
		logrus_logger.Log.Warn("handlers.go/CreateFriendshipHandler - error while decoding friendship: ", err)
		return
	}
	// добавляем запись о дружбе в бд
	if err := postgres.InsertFriendship(friendship.InvitedUserId, friendship.InvitingUserId); err != nil {
		http.Error(w, "handlers.go/CreateFriendshipHandler - error while creating friendship", http.StatusInternalServerError)
		logrus_logger.Log.Warn("handlers.go/CreateFriendshipHandler - error creating friendship: ", err)
		return
	}
	// уведомляем клиент об успешном завершении процесса
	w.WriteHeader(http.StatusCreated)
}

// GetUserTasksHandler возвращает список заданий пользователя по id
func GetUserTasksHandler(w http.ResponseWriter, r *http.Request) {
	logrus_logger.Log.Infof("received request %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
	// получаем id пользователя из запроса
	userId := chi.URLParam(r, "id")
	// парсим его в int64
	id, err := string_utils.ParseID(userId)
	if err != nil {
		http.Error(w, "handlers.go/GetUserTasksHandler - error while parsing user id", http.StatusBadRequest)
		logrus_logger.Log.Warn("handlers.go/GetUserTasksHandler - error while parsing user id: ", err)
		return
	}
	// получаем задания из бд
	tasks, err := postgres.SelectUserTasks(id)
	if err != nil {
		http.Error(w, "handlers.go/GetUserTasksHandler - error while selecting tasks", http.StatusInternalServerError)
		logrus_logger.Log.Warn("handlers.go/GetUserTasksHandler - error while selecting tasks: ", err)
		return
	}
	// энкодим данные о заданиях
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(tasks); err != nil {
		http.Error(w, "handlers.go/GetUserTasksHandler - error while selecting tasks", http.StatusInternalServerError)
		logrus_logger.Log.Warn("handlers.go/GetUserTasksHandler - error while selecting tasks: ", err)
		return
	}
}

// CreateTaskHandler создаёт новое задание, принимая информацию о нём в json формате
func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	task := &models.Task{}
	// декодим информацию о задании
	err := json.NewDecoder(r.Body).Decode(task)
	if err != nil {
		http.Error(w, "handlers.go/CreateTaskHandler - error while decoding task", http.StatusBadRequest)
		logrus_logger.Log.Warn("handlers.go/CreateTaskHandler - error while decoding task: ", err)
		return
	}
	// добавляем полученное задание в бд
	err = postgres.InsertTask(task)
	if err != nil {
		http.Error(w, "handlers.go/CreateTaskHandler - error while inserting task", http.StatusInternalServerError)
		logrus_logger.Log.Warn("handlers.go/CreateTaskHandler - error while inserting task: ", err)
		return
	}
	// отправляем клиенту статус успешного создания задания
	w.WriteHeader(http.StatusCreated)
}
