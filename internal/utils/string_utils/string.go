package string_utils

import "strconv"

// StringToInt64 парсит строку в int64
//
// Входным параметром функции является строка,
// которую нужно распарсить в int64.
//
// Выходными параметрами функции являются
// число, которое получилось в результате парсинга
// строки, и ошибка, если она возникла, в противном
// случае вместо неё будет возврашён nil.
func StringToInt64(str string) (int64, error) {
	id, err := strconv.ParseInt(str, 10, 0)
	if err != nil {
		return 0, err
	}
	return id, nil
}
