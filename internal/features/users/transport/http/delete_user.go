package users_transport_http

import (
	"net/http"

	core_logger "github.com/fffeehh/practice_backend_todoapp/internal/core/logger"
	core_http_response "github.com/fffeehh/practice_backend_todoapp/internal/core/transport/http/response"
	core_http_request "github.com/fffeehh/practice_backend_todoapp/internal/core/transport/http/request"
)

func(h *UsersHTTPHandler) DeleteUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err, 
			"failed to get userID path value",
			)
		return
	}

	if err := h.usersService.DeleteUser(ctx, userID); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to delete user",
			)
		return
	}
	responseHandler.NoContentResponse()
}
