package user_transport_http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/krxzy37/todoapp-go/internal/core/domain"
	core_logger "github.com/krxzy37/todoapp-go/internal/core/logger"
	core_http_request "github.com/krxzy37/todoapp-go/internal/core/transport/http/request"
	core_http_response "github.com/krxzy37/todoapp-go/internal/core/transport/http/response"
	core_http_types "github.com/krxzy37/todoapp-go/internal/core/transport/http/types"
	core_http_utils "github.com/krxzy37/todoapp-go/internal/core/transport/http/utils"
)

type PathUserRequest struct {
	FullName    core_http_types.Nullable[string] `json:"full_name"`
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number"`
}

func (r *PathUserRequest) Validate() error {
	if r.FullName.Set {
		if r.FullName.Value == nil {
			return fmt.Errorf("FullName cant be NULL")
		}

		FullNameLen := len([]rune(*r.FullName.Value))

		if FullNameLen < 3 || FullNameLen > 100 {
			return fmt.Errorf("'FullName' must be between 3 and 100 symbols")
		}
	}

	if r.PhoneNumber.Set {
		if r.PhoneNumber.Value != nil {
			phoneNumberLen := len([]rune(*r.PhoneNumber.Value))

			if phoneNumberLen < 10 || phoneNumberLen > 15 {
				return fmt.Errorf(
					"'PhoneNumber' must be betweeb 10 and 15 symbols",
				)

			}
			if !strings.HasPrefix(*r.PhoneNumber.Value, "+") {
				return fmt.Errorf("'PhoneNumber' must startswith '+'")
			}
		}
	}
	return nil
}

type PatchUserResponse UserDtoResponse

func (h *UsersHTTPHandler) PatchUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponceHandler(log, rw)

	userID, err := core_http_utils.GetIntPathParam(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get userID path value",
		)
	}

	var request PathUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(

			err,
			"failed to decode and validate http request",
		)
		return
	}

	userPatch := userPatchFromRequest(request)

	userDomain, err := h.usersService.PatchUser(ctx, userID, userPatch)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to patch user",
		)
		return
	}

	response := PatchUserResponse(userDTOFromDomain(userDomain))

	responseHandler.JSONResponse(response, http.StatusOK)

	log.Debug(
		fmt.Sprintf(
			"PathuserRequest fields:\nFullName:'%v'\nPhoneNumber:'%v'",
			request.FullName,
			request.PhoneNumber,
		),
	)
	rw.WriteHeader(http.StatusOK)

}

func userPatchFromRequest(request PathUserRequest) domain.UserPatch {
	return domain.UserPatch{
		FullName:    request.FullName.ToDomain(),
		PhoneNumber: request.PhoneNumber.ToDomain(),
	}
}
