package string_utils

import "strconv"

// ParseID парсит строку с id пользователя в int64
func ParseID(userID string) (int64, error) {
	id, err := strconv.ParseInt(userID, 10, 0)
	if err != nil {
		return 0, err
	}
	return id, nil
}
