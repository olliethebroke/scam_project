package models

// Leader содержит информацию для отображения лидера лиги
type Leader struct {
	Id        int64  `json:"id"`
	Position  int16  `json:"position"`
	Firstname string `json:"firstname"`
	Blocks    int64  `json:"blocks"`
}

// CreateUserRequest содержит поле запроса от клиента на создание пользователя
type CreateUserRequest struct {
	Username string `json:"username"`
}

type SelectUserResponse struct {
	Id             int64  `json:"id"`
	Firstname      string `json:"firstname"`
	Blocks         int64  `json:"blocks"`
	Record         int    `json:"record"`
	DaysStreak     int    `json:"days_streak"`
	InvitedFriends int    `json:"invited_friends"`
	IsPremium      bool   `json:"is_premium"`
	League         int16  `json:"league"`
	Award          bool   `json:"award"`
}

// UpdateUserRequest содержит поля запроса от клиента на обновление данных о пользователе
type UpdateUserRequest struct {
	Blocks         int64 `json:"blocks"`
	Record         int   `json:"record"`
	InvitedFriends int   `json:"invited_friends"`
	IsPremium      bool  `json:"is_premium"`
}

// Friendship содержит id двух пользователей - приглашённого и приглащающего
type Friendship struct {
	InvitedUserId  int64 `json:"invited_user_id"`
	InvitingUserId int64 `json:"inviting_user_id"`
}

// Task содержит информацию о задании
type Task struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
	Reward      int    `json:"reward"`
	ActionType  string `json:"action_type"`
	ActionData  string `json:"action_data"`
}

// GetUserTaskResponse содержит ответ на запро
type GetUserTaskResponse struct {
	Task        Task `json:"task"`
	IsCompleted bool `json:"is_completed"`
}
