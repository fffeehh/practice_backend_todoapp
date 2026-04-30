package core_http_request

import (
	"encoding/json"
	"fmt"
	"net/http"

	core_errors "github.com/fffeehh/practice_backend_todoapp/internal/core/errors"
	"github.com/go-playground/validator"
)

var requestValidator = validator.New()

type validatable interface {
	Validate() error
}

// читаем, декодируем и валидируем запрос 
/*
Вычитываем из тела входящего запроса json и декодируем его в в переданную DTO (dest). 
Затем смотрим, есть ли у этой DTO кастомные правила валидации. Если есть, то применяем их. Если нет, валидируем как обычно 
*/
func DecodeAndValidateRequest(r *http.Request, dest any) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		// возвращаемую из Decode() ошибку подставляем вместо %v для контекста. Ошибка преобразовывается в обычную строку
		// А вместо %w врапаем ошибку из нашего пакета, чтобы handler.go смог определить статус код ошибки
		return fmt.Errorf(
			"decode json: %v: %w",
			err,
			core_errors.ErrInvalidArgument,
			)
	}

	// переменная, куда мы будем сохранять ошибку и потом ее обрабатывать, 
	// чтобы избавиться от лишнего кода, тк в if-else обработка одинаковая
	var (
		err error
	)

	//  проверяем, подходит ли переданная структура dest под интерфейс validatable
	// благодаря этой строчке мы узнаем, подходит ли переданный интерфейс any с названием dest под интерфейс validatable
	// проще говоря, мы узнаем, есть ли у той структуры, которая передается на место dest метод validate().
	// Если есть, то ok = True, иначе False
	v, ok := dest.(validatable)

	if ok {
		err = v.Validate()
	} else {
		err = requestValidator.Struct(dest)
	}

	if err != nil {
		return fmt.Errorf(
			"request validation: %v: %w", 
			err,
			core_errors.ErrInvalidArgument,
		)
	}

	return nil
}
