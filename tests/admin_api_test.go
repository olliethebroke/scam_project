package tests

import (
	"bytes"
	"context"
	modelAPI "crypto_scam/internal/api/handler/admin/model"
	"crypto_scam/internal/app"
	"encoding/json"
	"github.com/brianvoe/gofakeit/v7"
	"net/http"
	"net/http/httptest"
)

// TestCreateTaskHandler тестирует хендлер
// создания нового задания.
func (s *ApiSuite) TestCreateTaskHandler() {
	r := s.Require()

	// формируем структуру-тело запроса
	createTaskRequest := &modelAPI.CreateTaskRequest{
		Description: gofakeit.Breakfast(),
		Reward:      gofakeit.Int32(),
		ActionType:  gofakeit.MovieName(),
		ActionData:  gofakeit.Gender(),
	}

	// сериализуем структуру в json
	reqBody, err := json.Marshal(createTaskRequest)
	r.NoError(err, "failed to serialize json")

	// создаём запрос
	req, err := http.NewRequest("POST", app.AdminCreateTaskPostfix, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "tma "+adminInitDataRaw)
	r.NoError(err, "failed to create request")

	// создаём и инициализируем переменную-имплементацию интерфейса ResponseWriter
	resp := httptest.NewRecorder()

	// создаём и инициализируем указатель на маршрутизатор
	router := s.a.Router(context.Background())

	// маршрутизатор должен обработать запрос
	router.ServeHTTP(resp, req)

	// проверяем возвращаемый статус код
	r.Equal(http.StatusCreated, resp.Result().StatusCode)

	// парсим id задания, которое было получено
	// в результате его создания
	createTaskResponse := &modelAPI.CreateTaskResponse{}
	_ = json.NewDecoder(resp.Body).Decode(createTaskResponse)

	// делаем запрос в бд, чтобы считать созданное задание
	// в ответе не должно быть ошибки
	createdTask, err := s.a.ServiceProvider().DB(context.Background()).SelectTask(createTaskResponse.Id)
	r.NoError(err)

	// проверяем поля, отправленные в бд
	// и поля, полученные из бд на соответствие
	r.Equal(createTaskRequest.Description, createdTask.Description)
	r.Equal(createTaskRequest.Reward, createdTask.Reward)
	r.Equal(createTaskRequest.ActionType, createdTask.ActionType)
	r.Equal(createTaskRequest.ActionData, createdTask.ActionData)

	// добавляем созданное задание в слайс заданий
	tasks = append(tasks, createdTask)

}

// TestDeleteTaskHandler тестирует хендлер
// удаления задания по id.
func (s *ApiSuite) TestDeleteTaskHandler() {
	r := s.Require()

	// считываем все задания из бд
	existingTasks, err := s.a.ServiceProvider().DB(context.Background()).SelectTasks()
	r.NoError(err)

	idToDelete := existingTasks[0].Id

	// формируем структуру-тело запроса
	deleteTaskRequest := modelAPI.DeleteTaskRequest{
		// будем удалять первое из заданий
		Id: idToDelete,
	}

	// сериализуем структуру в json
	reqBody, _ := json.Marshal(deleteTaskRequest)

	// создаём запрос
	req, _ := http.NewRequest("DELETE", app.AdminDeleteTaskPostfix, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "tma "+adminInitDataRaw)

	// создаём и инициализируем переменную-имплементацию интерфейса ResponseWriter
	resp := httptest.NewRecorder()

	// создаём и инициализируем указатель на маршрутизатор
	router := s.a.Router(context.Background())

	// маршрутизатор должен обработать запрос
	router.ServeHTTP(resp, req)

	// проверяем возвращаемый статус код
	r.Equal(http.StatusOK, resp.Result().StatusCode)

	// считываем все задания из бд
	existingTasks, err = s.a.ServiceProvider().DB(context.Background()).SelectTasks()
	r.NoError(err)

	// если удалённое задание осталось в бд,
	// тест провален
	for _, task := range existingTasks {
		if task.Id == idToDelete {
			r.FailNow("deleted task exists")
		}
	}

	// передаём в слайс заданий обновлённый список
	tasks = existingTasks
}

// TestDeleteUserHandler тестирует хендлер
// удаления пользователя по id.
func (s *ApiSuite) TestDeleteUserHandler() {
	r := s.Require()

	// удаляем первого пользователя
	idToDelete := users[0].Id

	// проверяем наличие пользователя в бд
	// в ответе не должно быть ошибки
	_, err := s.a.ServiceProvider().DB(context.Background()).SelectUser(idToDelete)
	r.NoError(err)

	// формируем структуру-тело запроса
	deleteUserRequest := modelAPI.DeleteUserRequest{
		Id: idToDelete,
	}

	// сериализуем структуру в json
	reqBody, _ := json.Marshal(deleteUserRequest)

	// создаём запрос
	req, _ := http.NewRequest("DELETE", app.AdminDeleteUserPostfix, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "tma "+adminInitDataRaw)

	// создаём и инициализируем переменную-имплементацию интерфейса ResponseWriter
	resp := httptest.NewRecorder()

	// создаём и инициализируем указатель на маршрутизатор
	router := s.a.Router(context.Background())

	// маршрутизатор должен обработать запрос
	router.ServeHTTP(resp, req)

	// проверяем возвращаемый статус код
	r.Equal(http.StatusOK, resp.Result().StatusCode)

	// делаем запрос в бд, чтобы считать удалённого пользователя
	// в ответе должна быть ошибка
	_, err = s.a.ServiceProvider().DB(context.Background()).SelectUser(deleteUserRequest.Id)
	r.Error(err)

	// удаляем первого пользователя из слайса
	users = users[1:]
}

// TestGetTasksHandler тестирует хендлер
// получения всех заданий.
func (s *ApiSuite) TestGetTasksHandler() {
	r := s.Require()

	// создаём запрос
	req, _ := http.NewRequest("GET", app.AdminGetTasksPostfix, nil)
	req.Header.Set("Authorization", "tma "+adminInitDataRaw)

	// создаём и инициализируем переменную-имплементацию интерфейса ResponseWriter
	resp := httptest.NewRecorder()

	// создаём и инициализируем указатель на маршрутизатор
	router := s.a.Router(context.Background())

	// маршрутизатор должен обработать запрос
	router.ServeHTTP(resp, req)

	// проверяем возвращаемый статус код
	r.Equal(http.StatusOK, resp.Result().StatusCode)

	// декодируем данные из ответа обработчика
	responseTasks := make([]*modelAPI.GetTasksResponse, 0, 4)
	err := json.NewDecoder(resp.Body).Decode(&responseTasks)
	r.NoError(err)

	// проверяем, равна ли длина слайса заданий
	// с количеством заданий в бд
	n := len(tasks)
	if len(responseTasks) != n {
		r.FailNow("task slice length is not equal to ", n)
	}
}
