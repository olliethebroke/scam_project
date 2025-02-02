package model

// Leader описывает информацию о лидере.
//
// Id - идентификатор лидера,
// Position - позиция в рейтинге,
// Firstname - имя лидера,
// Blocks - количество блоков.
type Leader struct {
	Id        int64
	Position  int16
	Firstname string
	Blocks    int64
}
