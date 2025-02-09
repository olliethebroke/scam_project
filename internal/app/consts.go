package app

const (
	// роли
	creator = 2
	admin   = 1
	user    = 0

	createUserPostfix       = "/user/create/{id}"
	createFriendshipPostfix = "/friendship/create"

	// администраторские запросы
	AdminCreateTaskPostfix = "/task/create"
	AdminDeleteTaskPostfix = "/task/delete"
	AdminGetTasksPostfix   = "/tasks/get"
	AdminDeleteUserPostfix = "/user/delete"

	// пользовательские запросы
	GetUserPostfix        = "/user/get"
	UpdateUserPostfix     = "/user/update"
	GetLeaderboardPostfix = "/leaderboard/get"
	GetUserTasksPostfix   = "/user/tasks/get"
)
