package users_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/fffeehh/practice_backend_todoapp/internal/core/logger"
	core_http_response "github.com/fffeehh/practice_backend_todoapp/internal/core/transport/http/response"
	core_http_request "github.com/fffeehh/practice_backend_todoapp/internal/core/transport/http/request"
)

// DTO которая представляет http ответ
type GetUsersResponse []UserDTOResponse


func (h *UsersHTTPHandler) GetUsers(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	limit, offset, err := getLimitOffsetQueryParams(r)
	if err != nil {
		// благодаря методу ErrorResponse структуры ResponseHandler обработка ошибок происходит в пару строк
		// В то время как он под копотом статус код, логгирует, формирует ответный json
		responseHandler.ErrorResponse(
			err,
			"failed to get 'limit'/'offset' query param",
			)
		return
	}

	userDomains, err := h.usersService.GetUsers(ctx, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get users",
			)
		return
	}

	// приводим ответ к нашему созданному виду DTO для чистоты и понятности
	response := GetUsersResponse(usersDTOFromDomains(userDomains))
	
	responseHandler.JSONResponse(response, http.StatusOK)
}


// получение query параметров limit и offset
func getLimitOffsetQueryParams(r *http.Request) (*int, *int, error) {
	const (
		limitQueryParamKey = "limit"
		offsetQueryParamKey = "offset"
	)

	limit, err := core_http_request.GetIntQueryParam(r, "limit")
	if err != nil {
		return nil, nil, fmt.Errorf("get 'limit' query param: %w", err)
	}

	offset, err := core_http_request.GetIntQueryParam(r, "offset")
	if err != nil {
		return nil, nil, fmt.Errorf("get 'offset' query param: %w", err)
	}

	return limit, offset, nil
}
