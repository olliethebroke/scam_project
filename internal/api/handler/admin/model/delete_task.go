package model

// DeleteTaskRequest описывает задание,
// которое нужно удалить.
//
// Id - идентификатор задания.
type DeleteTaskRequest struct {
	Id int16 `json:"id"`
}
