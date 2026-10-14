package user_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/krxzy37/todoapp-go/internal/core/logger"
	core_http_request "github.com/krxzy37/todoapp-go/internal/core/transport/http/request"
	core_http_response "github.com/krxzy37/todoapp-go/internal/core/transport/http/response"
)

type GetUsersResponse []UserDtoResponse

func (h *UsersHTTPHandler) GetUsers(rw http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponceHandler(log, rw)

	limit, offset, err := getLimitOffsetQuerryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get `limit/offset` query param")
		return
	}

	UserDomains, err := h.usersService.GetUsers(ctx, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get users")
	}

	response := GetUsersResponse(usersDTOFromDomain(UserDomains))

	responseHandler.JSONResponse(response, http.StatusOK)

}

func getLimitOffsetQuerryParams(r *http.Request) (*int, *int, error) {

	const (
		limitQueryParamKey  = "limit"
		offsetQueryParamKey = "offset"
	)

	limit, err := core_http_request.GetIntQuerryParam(r, limitQueryParamKey)
	if err != nil {
		return nil, nil, fmt.Errorf("get `limit` querry param: %w", err)
	}

	offset, err := core_http_request.GetIntQuerryParam(r, offsetQueryParamKey)
	if err != nil {
		return nil, nil, fmt.Errorf("get `offset` querry param: %w", err)
	}

	return limit, offset, nil

}
