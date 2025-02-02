package model

// CreateTaskRequest описывает данные
// о задании, которые пришли от клиента.
//
// Description - описание задания,
// Reward - награда за выполнение в блоках,
// ActionType - тип действия для выполнения,
// ActionData - подробности действия.
type CreateTaskRequest struct {
	Description string `json:"description"`
	Reward      int    `json:"reward"`
	ActionType  string `json:"action_type"`
	ActionData  string `json:"action_data"`
}
