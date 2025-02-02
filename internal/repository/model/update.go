package model

// Update структура описывающая данные
// для обновления игрового результата пользователя.
//
// Blocks - количество блоков,
// Record - игровой рекорд пользователя.
type Update struct {
	Blocks int64
	Record int64
}
