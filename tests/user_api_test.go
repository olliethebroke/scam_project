package tests

import (
	"bytes"
	"context"
	modelAPI "crypto_scam/internal/api/handler/user/model"
	"crypto_scam/internal/app"
	modelDB "crypto_scam/internal/repository/model"
	"encoding/json"
	"github.com/brianvoe/gofakeit/v7"
	tma "github.com/telegram-mini-apps/init-data-golang"
	"net/http"
	"net/http/httptest"
)

// TestGetLeaderboardHandler тестирует обработчик
// получения лидерборда.
func (s *ApiSuite) TestGetLeaderboardHandler() {
	r := s.Require()

	// формируем запрос
	req, err := http.NewRequest("GET", app.GetLeaderboardPostfix, nil)
	req.Header.Set("Authorization", "tma "+adminInitDataRaw)
	r.NoError(err)

	// создаём объект для чтения ответа обработчика,
	// реализующий интерфейс ResponseWriter
	resp := httptest.NewRecorder()

	// создаём и инициализируем указатель на маршрутизатор
	router := s.a.Router(context.Background())

	// обрабатываем запрос
	router.ServeHTTP(resp, req)

	// проверяем, что ответ положительный
	r.Equal(http.StatusOK, resp.Result().StatusCode)

	// декодируем ответ
	leadersResponse := make(map[int16][]*modelDB.Leader)
	err = json.NewDecoder(resp.Body).Decode(&leadersResponse)
	r.NoError(err)

	// сравнивем количество лидеров в четвёртой лиге
	// в ответе с данными из бд
	n := len(leaders[3])
	if len(leadersResponse[3]) != n {
		s.FailNow("leaders slice length is not equal to ", n)
	}
}

// TestGetUserHandler тестирует обработчик
// получения пользователя.
func (s *ApiSuite) TestGetUserHandler() {
	r := s.Require()

	// формируем запрос
	req, err := http.NewRequest("GET", app.GetUserPostfix, nil)
	req.Header.Set("Authorization", "tma "+adminInitDataRaw)
	r.NoError(err)

	// создаём объект для чтения ответа обработчика,
	// реализующий интерфейс ResponseWriter
	resp := httptest.NewRecorder()

	// создаём и инициализируем указатель на маршрутизатор
	router := s.a.Router(context.Background())

	// обрабатываем запрос
	router.ServeHTTP(resp, req)

	// проверяем, что ответ положительный
	r.Equal(http.StatusOK, resp.Result().StatusCode)

	// декодируем ответ
	getUserResponse := &modelAPI.GetUserResponse{}
	err = json.NewDecoder(resp.Body).Decode(getUserResponse)
	r.NoError(err)

	// парсим данные о пользователе, сделавшем запрос
	requestedUser, err := tma.Parse(adminInitDataRaw)
	r.NoError(err)

	// сравниваем поля ответа и пользователя,
	// сделавшего запрос
	r.Equal(getUserResponse.Id, requestedUser.User.ID)
	r.Equal(getUserResponse.Firstname, requestedUser.User.FirstName)

}

// TestGetUserTasksHandler тестирует обработчик
// получения заданий пользователя.
func (s *ApiSuite) TestGetUserTasksHandler() {
	r := s.Require()

	// формируем запрос
	req, err := http.NewRequest("GET", app.GetUserTasksPostfix, nil)
	req.Header.Set("Authorization", "tma "+adminInitDataRaw)
	r.NoError(err)

	// создаём объект для чтения ответа обработчика,
	// реализующий интерфейс ResponseWriter
	resp := httptest.NewRecorder()

	// создаём и инциализирует указатель на маршутизатор
	router := s.a.Router(context.Background())

	// обрабатываем запрос
	router.ServeHTTP(resp, req)

	// проверяем, что ответ положительный
	r.Equal(http.StatusOK, resp.Result().StatusCode)

	// декодируем ответ
	var getUserTasksResponse []*modelAPI.GetUserTasksResponse
	err = json.NewDecoder(resp.Body).Decode(&getUserTasksResponse)
	r.NoError(err)

	// сравниваем задания ответа с заданиями
	// из базы данных
	for _, userTask := range getUserTasksResponse {
		for _, completedTask := range completedTasks[adminId] {
			if userTask.Id == completedTask {
				if !userTask.IsCompleted {
					r.FailNow("not completed task is in 'completed' slice")
				}
			}
		}
	}

}

// TestUpdateUserHandler тестирует обработчик
// обновления игровых показателей пользователя.
func (s *ApiSuite) TestUpdateUserHandler() {
	r := s.Require()

	// формируем структуру-тело запрос
	updateUserRequest := &modelAPI.UpdateUserRequest{
		Blocks: gofakeit.Int64(),
		Record: gofakeit.Int64(),
	}

	// сериализуем структуру в json
	reqBody, err := json.Marshal(updateUserRequest)
	r.NoError(err)

	// формируем запрос
	req, err := http.NewRequest("PATCH", app.UpdateUserPostfix, bytes.NewBuffer(reqBody))
	req.Header.Set("Authorization", "tma "+adminInitDataRaw)
	r.NoError(err)

	// создаём объект для чтения ответа обработчика,
	// реализующий интерфейс ResponseWriter
	resp := httptest.NewRecorder()

	// создаём и инициализируем указатель на маршрутизатор
	router := s.a.Router(context.Background())

	// обрабатываем запрос
	router.ServeHTTP(resp, req)

	// проверяем, что ответ положительный
	r.Equal(http.StatusOK, resp.Result().StatusCode)

	// получаем пользователя из бд
	updatedUser, err := s.a.ServiceProvider().DB(context.Background()).SelectUser(adminId)
	r.NoError(err)

	// сравниваем поля обновленного пользователя
	// с теми, которые передавали в обработчик
	r.Equal(updateUserRequest.Record, updatedUser.Record)
	r.Equal(updateUserRequest.Blocks, updatedUser.Blocks)
}
