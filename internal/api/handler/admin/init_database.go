package admin_api

import "crypto_scam/internal/repository"

var database repository.Repository

// InitDatabase инициализирует переменную database.
//
// Входным параметром функции является переменная db,
// реализующая интерфейс Repository.
func InitDatabase(db repository.Repository) {
	database = db
}
