package model

// User отображает данные пользователя.
//
// Id - идентификатор пользователя в телеграм,
// Firstname - имя пользователя в телеграм,
// Blocks - количество блоков пользователя,
// Record - рекорд пользователя,
// DaysStreak - серия ежедневных заходов,
// InvitedFriends - количество приглашённых пользователей,
// IsPremium - купил ли пользователь премиум статус,
// League - лига, в которой находится пользователь,
// Award - нужно ли выдать награду за вход.
type User struct {
	Id             int64
	Firstname      string
	Blocks         int64
	Record         int64
	DaysStreak     int16
	InvitedFriends int16
	IsPremium      bool
	League         int16
	Award          bool
}
