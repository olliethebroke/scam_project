package model

// GetTasksResponse описывает данные задания,
// которое было запрошено клиентом.
//
// Id - идентификатор задания,
// Description - описание задания,
// Reward - награда за выполнение в блоках,
// ActionType - тип действия для выполнения,
// ActionData - подробности действия.
type GetTasksResponse struct {
	Id          int16  `json:"id"`
	Description string `json:"description"`
	Reward      int32  `json:"reward"`
	ActionType  string `json:"action_type"`
	ActionData  string `json:"action_data"`
}
