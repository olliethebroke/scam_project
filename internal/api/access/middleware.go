package access

import (
	"context"
	"crypto_scam/internal/config"
	"crypto_scam/internal/logger"
	"crypto_scam/internal/repository"
	"errors"
	tma "github.com/telegram-mini-apps/init-data-golang"
	"net/http"
	"strings"
)

// UserAuthMiddleware проверяет пользователя, отправившего апи запрос.
// Проверка осуществляется по правам доступа к тому или иному хендлеру,
// используется initDataRaw для определения подлинности id пользователя.
//
// Входными параметрами функции являются необходимая роль для выполнения запроса,
// объект, реализующий интерфейс Repository, для работы с базой данных и
// объект, реализующий интерфейс TGConfig, для работы с Telegram-API.
//
// Выходным параметром функции является функция, реализующая интерфейс Handler.
func UserAuthMiddleware(requiredRole int16, db repository.Repository, tg config.TGConfig) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// добавим логгирование запроса
			logger.Infof("received request %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

			// обрабатываем данные пользователя,
			// сделавшего запрос
			initData, err := getInitData(r, tg)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// проверяем, есть ли пользователь в бд
			// получаем его id и роль в переменную userRole
			userRole, err := db.SelectUserRole(initData.User.ID)
			if err != nil {
				http.Error(w, "user not found", http.StatusNotFound)
				return
			}

			// проверяем, есть ли у пользователя права
			// на использование хендлера
			if userRole.Role >= requiredRole {
				// создаём контекст с реализацией интерфейса бд
				// и id пользователя, по которому делается запрос
				ctx := r.Context()
				ctx = context.WithValue(ctx, "id", userRole.Id)
				ctx = context.WithValue(ctx, "db", db)

				// добавляем созданный контекст в запрос
				r = r.WithContext(ctx)

				// передаём управление конечному обработчику
				next.ServeHTTP(w, r)
			} else {
				// если у пользователя нет прав на запрос,
				// то отдаём ошибку
				http.Error(w, "access forbidden", http.StatusForbidden)
				return
			}
		})
	}
}

// getInitData отвечает за обработку данных пользователя
// из заголовка запроса.
// Функция извлекает данные из заголовка, проверяет их на подлинность,
// а затем парсит, в конце отдавая струткуру InitData.
//
// Входными параметрами функции являются указатель на тип Request
// и объект, реализующий интерфейс TGConfig.
//
// Выходными параметрами функции являются указатель на тип InitData
// и ошибка, если она возникла, в противном случае вместо неё будет
// возвращён nil.
func getInitData(r *http.Request, tg config.TGConfig) (*tma.InitData, error) {
	// извлекаем данные из заголовка запроса
	initDataRaw := r.Header.Get("Authorization")
	if initDataRaw == "" {
		return nil, errors.New("authorization header is missing")
	}

	// разбиваем значение по ключу Authorization на две части
	splitToken := strings.Split(initDataRaw, "tma ")
	if len(splitToken) != 2 {
		return nil, errors.New("authorization header is missing")
	}

	// сохраняем ту, которая относится к initDataRaw,
	// предоставленной телеграмом
	initDataRaw = splitToken[1]

	// валидируем данные пользователя
	if err := tma.Validate(initDataRaw, tg.Token(), tg.InitDataExpiration()); err != nil {
		return nil, errors.New("authorization header is missing")
	}

	// если данные валидны, то парсим их
	initData, err := tma.Parse(initDataRaw)
	if err != nil {
		return nil, errors.New("failed to parse init data")
	}

	// возвращаем данные пользователя
	// в виде структуры InitData
	return &initData, nil
}
