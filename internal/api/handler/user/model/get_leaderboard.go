package model

// GetLeaderboardResponse содержит информацию для отображения лидера лиги.
//
// Id - идентификатор лидера в телеграм,
// Position - позиция в рейтинге,
// Firstname - имя лидера,
// Blocks - количество блоков.
type GetLeaderboardResponse struct {
	Id        int64  `json:"id"`
	Position  int16  `json:"position"`
	Firstname string `json:"firstname"`
	Blocks    int64  `json:"blocks"`
}
