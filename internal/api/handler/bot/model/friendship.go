package model

// Friendship содержит id двух пользователей - приглашённого и приглащающего
type Friendship struct {
	InvitedUserId  int64 `json:"invited_user_id"`
	InvitingUserId int64 `json:"inviting_user_id"`
}
