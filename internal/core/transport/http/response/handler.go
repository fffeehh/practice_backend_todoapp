package core_http_response

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	core_errors "github.com/fffeehh/practice_backend_todoapp/internal/core/errors"
	core_logger "github.com/fffeehh/practice_backend_todoapp/internal/core/logger"
	"go.uber.org/zap"
)


type HTTPResponseHandler struct {
	log *core_logger.Logger
	rw http.ResponseWriter
}

func NewHTTPResponseHandler(
	log *core_logger.Logger,
	rw http.ResponseWriter,
) *HTTPResponseHandler {
	return &HTTPResponseHandler{
		log: log,
		rw: rw,
	}
}


func (h *HTTPResponseHandler) JSONResponse(
	responseBody any,
	statusCode int,
) {
	// Отдаем статус код об ошибке сервера
	h.rw.WriteHeader(statusCode)

	if err := json.NewEncoder(h.rw).Encode(responseBody); err != nil {
		h.log.Error("write HTTP response", zap.Error(err))
	}
}

func (h *HTTPResponseHandler) NoContentResponse() {
	h.rw.WriteHeader(http.StatusNoContent)
}


// метод для отправки http ответа в случае возникновения ошибки
func (h *HTTPResponseHandler) ErrorResponse(err error, msg string) {
	var (
		statusCode int
		logFunc func(string, ...zap.Field)
	)

	switch {
	case errors.Is(err, core_errors.ErrInvalidArgument):
		statusCode = http.StatusBadRequest
		// выбираем уровень Warn, потому что это не наша ошибка, 
		// но при этом стоит перепроверить все связи с фронтом и узнать, почему летят невалидные данные.
		logFunc = h.log.Warn

	case errors.Is(err, core_errors.ErrNotFound):
		statusCode = http.StatusNotFound
		// не наша ошибка, а неверные действия пользователя
		logFunc = h.log.Debug

	case errors.Is(err, core_errors.ErrConflict):
		statusCode = http.StatusConflict
		logFunc = h.log.Warn

	default:
		statusCode = http.StatusInternalServerError
		logFunc = h.log.Error
	}

	logFunc(msg, zap.Error(err))

	h.errorResponse(
		statusCode,
		err,
		msg,
		)
}

// метод для отправки http ответа в случае возникновения паники
func (h *HTTPResponseHandler) PanicResponse(p any, msg string) {
	statusCode := http.StatusInternalServerError
	err := fmt.Errorf("unexpected panic %v", p)

	// Логгируем панику
	h.log.Error(msg, zap.Error(err))
	
	h.errorResponse(
		statusCode,
		err,
		msg,
		)
}


func (h *HTTPResponseHandler) errorResponse(
	statusCode int,
	err error,
	msg string,
) {
	// Отдаем json с информацией о произошедшем в теле ответа
	response := map[string]string{
		"message": msg,
		"error": err.Error(),
	}

	h.JSONResponse(
		response,
		statusCode,
		)
}
