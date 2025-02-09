package app

import (
	"context"
	"crypto_scam/internal/api/access"
	adminAPI "crypto_scam/internal/api/handler/admin"
	userAPI "crypto_scam/internal/api/handler/user"
	"crypto_scam/internal/closer"
	"crypto_scam/internal/config"
	"crypto_scam/internal/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// App описывает поля, необходимые
// для работы приложения.
//
// serviceProvider - структура, реализующая зависимости.
// r - маршрутизатор.
type App struct {
	serviceProvider *serviceProvider
	r               *chi.Mux
}

// Run запускает HTTP сервер, добавляя в конец завершения
// метода освобождение ресурсов.
//
// Выходным параметром метода является ошибка,
// если она возникла, в противном случае вместо
// неё будет возвращён nil.
func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	return a.runHTTPServer()
}

// NewApp создаёт новый объект типа App.
//
// Входным параметром функции ялвяется контекст,
//
// Выходным параметром функции являются указатель на тип App
// и ошибка, если она возникла, в противном случае вместо неё
// будет возвращён nil.
func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	// инициализируем зависимости
	if err := a.initDeps(ctx); err != nil {
		return nil, err
	}

	// возвращаем указатель на тип App
	// и nil в качестве ошибки
	return a, nil
}

// initDeps инициализирует зависимости приложения.
//
// Входным параметром метода является контекст.
//
// Выходным параметром метода является ошибка,
// если она возникла, в противном случае вместо
// неё будет возвращён nil.
func (a *App) initDeps(ctx context.Context) error {
	// создаём слайс функций инициализации зависимостей
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initHTTPServer,
	}

	// вызываем методы из слайса inits
	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	// если не произошло ошибок
	// возвращаем nil
	return nil
}

// initConfig загружает переменные из env файла в переменные
// окружения процесса.
//
// Входным параметром метода является контекст.
//
// Выходным параметром метода является ошибка,
// если она возникла, в противном случае вместо
// неё будет возвращён nil.
func (a *App) initConfig(_ context.Context) error {
	// получаем путь к файлу окружения с помощью
	// переменной окружения процесса CONFIG_PATH
	configPath := os.Getenv("CONFIG_PATH")
	if len(configPath) == 0 {
		configPath = filepath.Join(".", "..", "test.env")
		logger.Warn("failed to get CONFIG_PATH env variable; test.env is used")
	}
	logger.Info("path to config:", configPath)

	// проверяем, существует ли файл
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("config file does not exist at path: ", configPath)
	}

	// получаем абсолютный путь к файлу
	var err error
	configPath, err = filepath.Abs(configPath)
	if err != nil {
		log.Fatal("failed to get absolute path to env file")
	}

	// загружаем переменные из указанного .env файла в переменные окружения процесса
	if err = config.Load(configPath); err != nil {
		logger.Fatal("failed to load variables from env file: ", err)
	}

	// если ошибок не возникло,
	// возвращаем nil
	return nil
}

// initServiceProvider инициализиурет поле serviceProvider.
//
// Входным параметром метода является контекст.
//
// Выходным параметром метода является ошибка,
// если она возникла, в противном случае вместо
// неё будет возвращён nil.
func (a *App) initServiceProvider(_ context.Context) error {
	// инициализируем поле serviceProvider
	a.serviceProvider = newServiceProvider()
	return nil
}

// initHTTPServer инициализирует и настраивает маршрутизатор
// для работы сервера по HTTP протоколу.
//
// Входным параметром метода является контекст.
//
// Выходным параметром метода является ошибка,
// если она возникла, в противном случае вместо
// неё будет возвращён nil.
func (a *App) initHTTPServer(ctx context.Context) error {
	// инициализируем маршрутизатор
	a.r = chi.NewRouter()

	// если где-то внутри сервера (обработчика запроса) произойдет паника, приложение не должно упасть
	a.r.Use(middleware.Recoverer)

	// инициализируем структуру для работы с бд
	a.serviceProvider.DB(ctx)

	// инициализируем структуру для работы с тг
	a.serviceProvider.TGConfig()

	// администраторские запросы
	a.r.With(access.UserAuthMiddleware(admin, a.serviceProvider.db, a.serviceProvider.tgConfig)).
		Get(AdminGetTasksPostfix, adminAPI.GetTasksHandler)
	a.r.With(access.UserAuthMiddleware(admin, a.serviceProvider.db, a.serviceProvider.tgConfig)).
		Post(AdminCreateTaskPostfix, adminAPI.CreateTaskHandler)
	a.r.With(access.UserAuthMiddleware(admin, a.serviceProvider.db, a.serviceProvider.tgConfig)).
		Delete(AdminDeleteTaskPostfix, adminAPI.DeleteTaskHandler)
	a.r.With(access.UserAuthMiddleware(admin, a.serviceProvider.db, a.serviceProvider.tgConfig)).
		Delete(AdminDeleteUserPostfix, adminAPI.DeleteUserHandler)
	// пользовательские запросы
	a.r.With(access.UserAuthMiddleware(user, a.serviceProvider.db, a.serviceProvider.tgConfig)).
		Get(GetUserPostfix, userAPI.GetUserHandler)
	a.r.With(access.UserAuthMiddleware(user, a.serviceProvider.db, a.serviceProvider.tgConfig)).
		Get(GetLeaderboardPostfix, userAPI.GetLeaderboardHandler)
	a.r.With(access.UserAuthMiddleware(user, a.serviceProvider.db, a.serviceProvider.tgConfig)).
		Get(GetUserTasksPostfix, userAPI.GetUserTasksHandler)
	a.r.With(access.UserAuthMiddleware(user, a.serviceProvider.db, a.serviceProvider.tgConfig)).
		Patch(UpdateUserPostfix, userAPI.UpdateUserHandler)

	return nil
}

// runHTTPServer запускает http сервер.
//
// Выходным параметром метода является ошибка,
// если она возникла, в противном случае вместо
// неё будет возвращён nil.
func (a *App) runHTTPServer() error {
	// запускаем сервер
	logger.Info("server is listening")
	err := http.ListenAndServe(a.serviceProvider.HTTPConfig().Address(), a.r)
	if err != nil {
		return err
	}
	return nil
}

// ServiceProvider возвращает указатель на структуру зависимостей.
//
// Выходным параметром метода является указатель на тип serviceProvider.
func (a *App) ServiceProvider() *serviceProvider {
	// если указатель serviceProvider nil
	if a.serviceProvider == nil {
		// инициализируем поле serviceProvider
		a.serviceProvider = newServiceProvider()
	}

	// возвращаем указатель на serviceProvider
	return a.serviceProvider
}

// Router возвращает указатель на маршрутизатор.
//
// Выходным параметром метода является указатель на тип chi.Mux.
func (a *App) Router(ctx context.Context) *chi.Mux {
	// если указатель на маршрутизатор nil
	if a.r == nil {
		// инициализируем маршрутизатор
		_ = a.initHTTPServer(ctx)
	}

	// возвращаем указатель на chi.Mux
	return a.r
}
