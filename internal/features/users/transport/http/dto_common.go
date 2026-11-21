package user_transport_http

import "github.com/krxzy37/todoapp-go/internal/core/domain"

type UserDtoResponse struct {
	ID          int     `json:"id" example:"10"`
	Version     int     `json:"version" example:"3"`
	FullName    string  `json:"full_name" example:"Ivan Petrenko"`
	PhoneNumber *string `json:"phone_number" example:"+38050734569"`
}

func userDTOFromDomain(user domain.User) UserDtoResponse {
	return UserDtoResponse{
		ID:          user.ID,
		Version:     user.Version,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
	}
}

func usersDTOFromDomain(users []domain.User) []UserDtoResponse {
	usersDTO := make([]UserDtoResponse, len(users))

	for i, user := range users {
		usersDTO[i] = userDTOFromDomain(user)
	}

	return usersDTO

}
