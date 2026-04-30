package core_http_utils

import (
	"fmt"
	"net/http"
	"strconv"

	core_errors "github.com/fffeehh/practice_backend_todoapp/internal/core/errors"
)

// пакет для инструментов по получению query параметров (например для get_users.go)

// функция для получения параметра типа int
func GetIntQueryParam(r *http.Request, key string) (*int, error) {
	// query параметр всегда приходит в виде строки. Если он не задан, значит приходит пустая строка
	param := r.URL.Query().Get(key)
	if param == "" {
		// если нам не передали параметр, то это не ошибка. 
		// Поэтому возвращаем из функции указатель на int и если параметра нет, возвращаем nil
		return nil, nil
	}

	val, err := strconv.Atoi(param) 
	if err != nil {
		return nil, fmt.Errorf(
			"param='%s' by key='%s' not a valid integer: %v: %w",
			param,
			key,
			err,
			core_errors.ErrInvalidArgument,
			)
	}

	return &val, nil
}
