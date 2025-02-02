package access

import "crypto_scam/internal/repository"

// InitDatabase инициализирует переменную database.
//
// Входным параметром функции является переменная db,
// реализующая интерфейс Repository.
func InitDatabase(db repository.Repository) {
	database = db
}
