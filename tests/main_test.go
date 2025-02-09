package tests

import (
	"context"
	"crypto_scam/internal/app"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type ApiSuite struct {
	suite.Suite
	a *app.App
}

// TestApiSuite отвечает за запуск тестовых методов,
// объединённых в сьют.
// Если запуск тестов вызван с флагом --short, то
// TestApiSuite будет пропущена, как и тесты из сьюта.
// Внутри suite.Run() фреймворк берет на себя управление,
// вызывая SetupSuite перед всеми тестами сьюта,
// затем тестовые методы, и в конце TearDownSuite.
func TestApiSuite(t *testing.T) {
	// если флаг --short
	if testing.Short() {
		// пропускаем тест
		t.Skip()
	}
	// запускаем тестовые методы струтуры ApiSuite
	suite.Run(t, new(ApiSuite))
}

// TestMain запускат тесты
func TestMain(m *testing.M) {
	rc := m.Run()
	os.Exit(rc)
}

// SetupSuite подготавливает suite к тестам.
// Происходит инициализация приложения и
// заполнение базы данных тестовыми данными.
func (s *ApiSuite) SetupSuite() {
	s.initApp()
	s.populateDB()
}

// initApp инициализирует зависимости приложения.
// Устанавливает соединение с базой данных.
func (s *ApiSuite) initApp() {
	var err error
	s.a, err = app.NewApp(context.Background())
	if err != nil {
		s.FailNow("failed to initialize the application")
	}
}

// populateDB заполняет тестовую базу данных тестовыми данными.
func (s *ApiSuite) populateDB() {
	// добавляем пользователей в бд
	for _, user := range users {
		if _, err := s.a.ServiceProvider().DB(context.Background()).InsertUser(user.Id, user.Firstname); err != nil {
			s.FailNow("failed to insert user", err)
		}
	}

	// добавляем задания
	for _, task := range tasks {
		if _, err := s.a.ServiceProvider().DB(context.Background()).InsertTask(task); err != nil {
			s.FailNow("failed to insert task", err)
		}
	}

	// добавляем выполненные задания
	for _, completedTask := range completedTasks[adminId] {
		if err := s.a.ServiceProvider().DB(context.Background()).
			InsertCompletedTask(adminId, completedTask); err != nil {
			s.FailNow("failed to insert completed task", err)
		}
	}

	// добавляем лидеров
	for _, leader := range leaders[3] {
		if err := s.a.ServiceProvider().DB(context.Background()).InsertLeader(leader, 3); err != nil {
			s.FailNow("failed to insert leader", err)
		}
	}

	// добавляем администратора
	if err := s.a.ServiceProvider().DB(context.Background()).InsertAdmin(adminId); err != nil {
		s.FailNow("failed to insert admin", err)
	}
}
