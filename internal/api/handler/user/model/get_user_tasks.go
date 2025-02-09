package model

// GetUserTasksResponse отображает данные о задании
// пользователя для отправки на клиент.
//
// Id - идентификатор пользователя,
// Description - описание задания,
// Reward - награда за выполнение задания,
// ActionType - тип действия для выполнения,
// ActionData - требуемое действие для выполнения,
// IsCompleted - выполнено ли задание.
type GetUserTasksResponse struct {
	Id          int16  `json:"id"`
	Description string `json:"description"`
	Reward      int32  `json:"reward"`
	ActionType  string `json:"action_type"`
	ActionData  string `json:"action_data"`
	IsCompleted bool   `json:"is_completed"`
}
