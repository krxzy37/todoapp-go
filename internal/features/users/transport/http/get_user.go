package user_transport_http

import (
	"net/http"

	core_logger "github.com/krxzy37/todoapp-go/internal/core/logger"
	core_http_response "github.com/krxzy37/todoapp-go/internal/core/transport/http/response"
	core_http_utils "github.com/krxzy37/todoapp-go/internal/core/transport/http/utils"
)

type GetUserResponse UserDtoResponse

func (h *UsersHTTPHandler) GetUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponceHandler(log, rw)

	userID, err := core_http_utils.GetIntPathParam(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get userID pathValue",
		)
		return
	}

	user, err := h.usersService.GetUser(ctx, userID)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get user",
		)
		return
	}

	response := GetUserResponse(userDTOFromDomain(user))

	responseHandler.JSONResponse(response, http.StatusOK)

}
