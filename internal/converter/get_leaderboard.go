package converter

import (
	modelAPI "crypto_scam/internal/api/handler/user/model"
	modelDB "crypto_scam/internal/repository/model"
)

// LeadersToGetLeaderboardResponse конвертирует мапу лидеров по лигам
// из типа данных modelDB.Leader в тип данных для отправки на клиент modelAPI.GetLeaderboardResponse.
//
// Входными данными для функции является мапа, где ключ - это идентификатор лиги (int16),
// а значение - срез указателей на модели лидеров из БД (modelDB.Leader).
//
// На выходе функция возвращает новую мапу, где ключи те же самые (идентификаторы лиг),
// а значения - срезы указателей на объекты типа modelAPI.GetLeaderboardResponse,
// которые содержат только необходимые данные для отображения на клиенте:
// Id, Blocks, Position и Firstname лидеров.
func LeadersToGetLeaderboardResponse(from map[int16][]*modelDB.Leader) map[int16][]*modelAPI.GetLeaderboardResponse {
	// создаём мапу
	// ключ - лига
	// значение - слайс указателей на GetLeaderboardResponse
	to := make(map[int16][]*modelAPI.GetLeaderboardResponse)

	// проходимся по списку лидеров в каждой лиге
	// и конвертируем данные в те, которые должны быть
	// отправлены на клиент
	for league, leaders := range from {
		newSlice := make([]*modelAPI.GetLeaderboardResponse, 0, len(leaders))
		for _, leader := range leaders {
			newLeader := &modelAPI.GetLeaderboardResponse{
				Id:        leader.Id,
				Blocks:    leader.Blocks,
				Position:  leader.Position,
				Firstname: leader.Firstname,
			}
			newSlice = append(newSlice, newLeader)
		}
		to[league] = newSlice
	}

	// возвращаем сконвертированную мапу
	return to
}
