package users_postgres_repository

import (
	"context"
	"fmt"

	"github.com/fffeehh/practice_backend_todoapp/internal/core/domain"
)

func (r *UsersRepository) GetUsers(
	ctx context.Context,
	limit *int,
	offset *int,
) ([]domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	// если limit или offset равен nil, тогда при подстановки через pgx значений в $1 и $2, они не будут учитываться. 
	// Работает именно с LIMIT и OFFSET
	query := `
	SELECT id, version, full_name, phone_number
	FROM todoapp.users
	ORDER BY id ASC
	LIMIT $1
	OFFSET $2
	`

	rows, err := r.pool.Query(
		ctx,
		query,
		limit,
		offset,
		)
	if err != nil {
		return nil, fmt.Errorf("select users: %w", err)
	}
	defer rows.Close()

	// идем по каждой строчке и каждую строчку преобразовываем в userModel,
	// а потом массив юзер модеелй преобразуем в юзер домены и вернем

	var userModels []UserModel
	for rows.Next() {
		var userModel UserModel

		err := rows.Scan(
			&userModel.ID,
			&userModel.Version,
			&userModel.FullName,
			&userModel.PhoneNumber,
			)
		if err != nil {
			return nil, fmt.Errorf("scan users: %w", err)
		}

		userModels = append(userModels, userModel)

		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("next rows: %w", err)
		}
	}

	userDomains := userDomainsFromModels(userModels)

	return userDomains, nil

}

