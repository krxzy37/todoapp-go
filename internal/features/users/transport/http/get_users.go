package user_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/krxzy37/todoapp-go/internal/core/logger"
	core_http_request "github.com/krxzy37/todoapp-go/internal/core/transport/http/request"
	core_http_response "github.com/krxzy37/todoapp-go/internal/core/transport/http/response"
)

type GetUsersResponse []UserDtoResponse

// GetUsers godoc
// @Summary users list
// @Description users list with optional pagination
// @Tags users
// @Produce json
// @Param limit query int false "page size with users"
// @Param offset query int false "page offset with users"
// @Success 200 {object} GetUsersResponse "Success get users list"
// @Failure 400 {object} core_http_response.ErrorResponse "Bad Request"
// @Failure 500 {object} core_http_response.ErrorResponse "Internal Server Error"
// @Router /users [get]
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

	limit, err := core_http_request.GetIntQueryParam(r, limitQueryParamKey)
	if err != nil {
		return nil, nil, fmt.Errorf("get `limit` querry param: %w", err)
	}

	offset, err := core_http_request.GetIntQueryParam(r, offsetQueryParamKey)
	if err != nil {
		return nil, nil, fmt.Errorf("get `offset` querry param: %w", err)
	}

	return limit, offset, nil

}
