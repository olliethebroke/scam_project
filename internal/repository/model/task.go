package model

// Task содержит информацию о задании.
//
// Id - идентификатор задания,
// Description - описание задания,
// Reward - награда за выполнение в блоках,
// ActionType - тип действия для выполнения,
// ActionData - подробности действия.
type Task struct {
	Id          int16
	Description string
	Reward      int32
	ActionType  string
	ActionData  string
	IsCompleted bool
}
