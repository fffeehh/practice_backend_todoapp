package users_transport_http

import "github.com/fffeehh/practice_backend_todoapp/internal/core/domain"

// общая DTO чтобы потом мы могли ее переиспользовать в других хэндлерах(create_user.go, get_users.go)
type UserDTOResponse struct {
	ID int `json:"id"`
	Version int `json:"version"`
	FullName string `json:"full_name"`
	PhoneNumber *string `json:"phone_number"`
}

func userDTOFromDomain(user domain.User) UserDTOResponse {
	return UserDTOResponse{
		ID: user.ID,
		Version: user.Version,
		FullName: user.FullName,
		PhoneNumber: user.PhoneNumber,
	}
}

// функция, которое преобразует массив доменных сущностей в массив DTO
func usersDTOFromDomains(users []domain.User) []UserDTOResponse {
	usersDTO := make([]UserDTOResponse, len(users))

	for i, user := range users {
		usersDTO[i] = userDTOFromDomain(user)
	}

	return usersDTO
}
