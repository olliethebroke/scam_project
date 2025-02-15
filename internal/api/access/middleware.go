package access

import (
	"context"
	"crypto_scam/internal/api/ratelimiter"
	"crypto_scam/internal/config"
	"crypto_scam/internal/logger"
	"crypto_scam/internal/repository"
	"errors"
	tma "github.com/telegram-mini-apps/init-data-golang"
	"net/http"
	"strings"
	"time"
)

const secondsForRequest = 2

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

			// из существующего контекста создаём новый,
			// добавляя таймаут
			ctx, cancel := context.WithTimeout(r.Context(), secondsForRequest*time.Second)
			defer cancel()

			// аутентифицируем пользователя
			initData, err := authenticate(r, tg)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// проверяем, не превышен ли лимит запросов
			isExceeded, err := isRateLimitExceeded(ctx, db, initData.User.ID)
			if err != nil {
				logger.Warn(err.Error())
				http.Error(w, "failed to check rate limit", http.StatusInternalServerError)
				return
			}
			if isExceeded {
				http.Error(w, "too many requests", http.StatusTooManyRequests)
				return
			}

			// авторизируем пользователя
			err = authorize(ctx, db, initData.User.ID, requiredRole)
			if err != nil {
				http.Error(w, err.Error(), http.StatusMethodNotAllowed)
				return
			}

			// в контекст добавляем идентификатор пользователя
			// и объект, реализующий интерфейс работы с бд
			ctx = context.WithValue(ctx, "id", initData.User.ID)
			ctx = context.WithValue(ctx, "db", db)

			// добавляем созданный контекст в запрос
			r = r.WithContext(ctx)

			// создаём канал done для уведомления
			// о завершении работы обработчика
			done := make(chan struct{})

			// создаём горутину
			go func() {
				// передаём управление конечному обработчику
				next.ServeHTTP(w, r)

				// после окончания работы
				// закрываем канал done
				close(done)
			}()

			// проверяем, успел ли обработчик завершиться
			// до таймаута
			select {
			case <-ctx.Done():
				// если не успел, то возвращаем на клиент ошибку
				http.Error(w, "the server didn't respond in time", http.StatusGatewayTimeout)
				return
			case <-done:
				// если успел, то всё крутяк вообще
				// завершаем выполнение функции
				return
			}

		})
	}
}

// authenticate отвечает за обработку данных пользователя
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
func authenticate(r *http.Request, tg config.TGConfig) (*tma.InitData, error) {
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

// authorize проверяет, имеет ли пользователь, совершивший запрос
// права на его исполнение.
//
// Входными параметрами функции являются контекст,
// объект, реализующий интерфейс взаимодействия с базой данных,
// идентификатор пользователя, совершающего запрос, и
// необходимая роль для выполнения запроса.
//
// Выходным параметром функции является ошибка, если она возникла,
// в противном случае будет возвращён nil.
func authorize(ctx context.Context, db repository.Repository, id int64, requiredRole int16) error {
	// проверяем, есть ли пользователь в бд
	// получаем его id и роль в переменную userRole
	userRole, err := db.SelectUserRole(ctx, id)
	if err != nil {
		return errors.New("user not found")
	}

	// проверяем, есть ли у пользователя права
	// на использование хендлера
	if userRole.Role < requiredRole {
		// если у пользователя нет прав на запрос,
		// то отдаём ошибку
		return errors.New("access forbidden")
	}

	// если ошибки не возникло
	// возвращаем nil
	return nil
}

// isRateLimitExceeded проверяет, превышен ли лимит запросов от конкретного пользователя.
//
// Входными параметрами функции являются контекст,
// объект, реализующий интерфейс взаимодействия с базой данных
// и идентификатор пользователя, совершающего запрос.
//
// Выходными параметрами функции являются булевая переменная, обозначающая, превысил ли
// пользователь лимит запосов, и ошибка, если она возникла, в противном случае вместо неё будет
// возвращён nil.
func isRateLimitExceeded(ctx context.Context, db repository.Repository, id int64) (bool, error) {
	// проверяем лимит запросов по id
	status := ratelimiter.LimitRequest(id)
	switch status {
	case 1:
		// пользователь не превысил лимит
		// можно выполнить запрос
		return false, nil
	case 0:
		// пользователь превысил лимит
		// запрос запрещён
		// необходимо сделать паузу в запросах
		return true, nil
	case -1:
		// пользователь сильно превысил лимит
		// запрос запрещён
		// удаляем аккаунт
		return true, db.DeleteUser(ctx, id)
	default:
		// в ином случае
		return true, errors.New("failed to check rate limits")
	}
}
