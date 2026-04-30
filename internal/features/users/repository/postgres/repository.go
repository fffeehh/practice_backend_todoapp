package users_postgres_repository

import core_postgres_pool "github.com/fffeehh/practice_backend_todoapp/internal/core/repository/postgres/pool"

type UsersRepository struct {
	// Тут мы можем напрямую зависеть от подключения к базе данных pgx.Conn, но это плохая практика, тк 
	// при тестировании репозитория придется устанавливать подключение и поднимать базу данных, что в unit тестах неприемлимо.
	// Поэтому мы должны зависеть от интерфейса подключения. Его реализацию выносим в core, он будет общий для всех фичей
	// Описываем его в core/repository/postgres/pool/pool.go
	pool core_postgres_pool.Pool
}

func NewUsersRepository(
	pool core_postgres_pool.Pool,
) *UsersRepository {
	return &UsersRepository{
		pool: pool,
	}
}
