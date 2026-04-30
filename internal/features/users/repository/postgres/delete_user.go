package users_postgres_repository

import (
	"context"
	"fmt"

	core_errors "github.com/fffeehh/practice_backend_todoapp/internal/core/errors"
)


func(r *UsersRepository) DeleteUser(
	ctx context.Context,
	id int,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	DELETE from todoapp.users
	WHERE id=$1
	`
// Exec() возвращает commandTag, в котором зашита информация о действии
	cmdTag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}

	// эта строчка значит, что запрос прошел успешно, но ни на какие строки он не повлиял. То есть такого пользователя нет
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("user with id='%d': %w", id, core_errors.ErrNotFound)
	}

	return nil 
}
