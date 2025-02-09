package model

// GetUserResponse отображает данные о пользователе
// для отправки на клиент.
//
// Id - идентификатор пользователя в телеграм,
// Firstname - имя пользователя в телеграм,
// Blocks - количество блоков пользователя,
// Record - игровой рекорд пользователя,
// DaysStreak - серия ежедневных заходов,
// InvitedFriends - количество приглашённых пользователей,
// IsPremium - купил ли пользователь премиум статус,
// League - лига, в которой находится пользователь,
// Award - нужно ли выдать награду за вход.
type GetUserResponse struct {
	Id             int64  `json:"id"`
	Firstname      string `json:"firstname"`
	Blocks         int64  `json:"blocks"`
	Record         int64  `json:"record"`
	DaysStreak     int16  `json:"days_streak"`
	InvitedFriends int16  `json:"invited_friends"`
	IsPremium      bool   `json:"is_premium"`
	League         int16  `json:"league"`
	Award          bool   `json:"award"`
}
