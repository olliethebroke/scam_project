package app

const (
	// роли
	creator = 2
	admin   = 1
	user    = 0

	createUserPostfix       = "/user/create/{id}"
	createFriendshipPostfix = "/friendship/create"

	// администраторские запросы
	adminCreateTaskPostfix = "/task/create"
	adminDeleteTaskPostfix = "/task/delete"
	adminGetTasksPostfix   = "/tasks/get"
	adminDeleteUserPostfix = "/user/delete"

	// пользовательские запросы
	getUserPostfix        = "/user/get"
	updateUserPostfix     = "/user/update"
	getLeaderboardPostfix = "/leaderboard/get"
	getUserTasksPostfix   = "/user/tasks/get"
)
