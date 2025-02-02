package access

import (
	"context"
	"crypto_scam/internal/logger"
	"crypto_scam/internal/repository"
	"crypto_scam/pkg/hooks/telegram"
	"errors"
	tma "github.com/telegram-mini-apps/init-data-golang"
	"net/http"
	"strings"
)

var database repository.Repository

// UserAuthMiddleware проверяет пользователя, отправившего апи запрос.
// Проверка осуществляется по правам доступа к тому или иному хендлеру,
// используется initDataRaw для определения подлинности id пользователя
func UserAuthMiddleware(requiredRole int16) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// добавим логгирование запроса
			logger.Infof("received request %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

			// обрабатываем данные пользователя,
			// сделавшего запрос
			initData, err := getInitData(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// проверяем, есть ли пользователь в бд
			// получаем его id и роль в переменную userInfo
			userInfo, err := database.SelectUserRole(initData.User.ID)
			if err != nil {
				http.Error(w, "user not found", http.StatusNotFound)
				return
			}

			// проверяем, есть ли у пользователя права
			// на использование хендлера
			if userInfo.Role >= requiredRole {
				// создаём контекст с id пользователя,
				// по которому делается запрос
				ctx := context.WithValue(r.Context(), "id", userInfo.Id)

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
// а затем парсит, в конце отдавая струткуру InitData
func getInitData(r *http.Request) (*tma.InitData, error) {
	// извлекаем данные из заголовка запроса
	initDataRaw := r.Header.Get("Authorization")
	if initDataRaw == "" {
		// http.Error(w, "authorization header is missing", http.StatusUnauthorized)
		return nil, errors.New("authorization header is missing")
	}

	// разбиваем значение по ключу Authorization на две части
	splitToken := strings.Split(initDataRaw, "tma ")
	if len(splitToken) != 2 {
		// http.Error(w, "invalid authorization format", http.StatusUnauthorized)
		return nil, errors.New("authorization header is missing")
	}

	// сохраняем ту, которая относится к initDataRaw,
	// предоставленной телеграмом
	initDataRaw = splitToken[1]

	// валидируем данные пользователя
	if err := tma.Validate(initDataRaw, telegram.TgConfig.Token(), telegram.TgConfig.InitDataExpiration()); err != nil {
		// http.Error(w, "bad init data", http.StatusBadRequest)
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
