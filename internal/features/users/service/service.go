package users_service

import (
	"context"

	"github.com/fffeehh/practice_backend_todoapp/internal/core/domain"
)

type UsersService struct {
	// создаем поля репозитория, тк сервис от него зависит
	usersRepository UsersRepository
}

// уровень сервиса зависит от интерфейса репозитория, поэтому создаем интерфейс 
type UsersRepository interface {
	CreateUser(
		ctx context.Context,
		user domain.User,
	) (domain.User, error)

	GetUsers(
		ctx context.Context,
		limit *int,
		offset *int,
	) ([]domain.User, error)

	GetUser(
		ctx context.Context,
		id int,
	) (domain.User, error)
	
	DeleteUser(
		ctx context.Context,
		id int,
	) error

	PatchUser(
		ctx context.Context,
		id int,
		user domain.User,
	) (domain.User, error)
}

func NewUsersService(
	usersRepository UsersRepository, 
) *UsersService {
	return &UsersService{
		usersRepository: usersRepository,
	}
}

// usersService должен имплементировать (соотствовать) интерфейсу сервиса, который мы описали на транспортном уровне в transport.go
