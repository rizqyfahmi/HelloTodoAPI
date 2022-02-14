package repositories

import (
	"TodoAPI/application/auth/entities"
	"context"
	"database/sql"
)

type AuthRepository struct {
	db *sql.DB
}

func InitAuthRepository(db *sql.DB) AuthRepositoryProtocol {
	return &AuthRepository{
		db: db,
	}
}

func (repo *AuthRepository) Store(ctx context.Context, user *entities.User) error {

	query := "INSERT INTO user(id, username, password, updated_at, created_at) VALUES (?, ?, ?, ?, ?)"
	stmt, errPrepare := repo.db.PrepareContext(ctx, query)

	if errPrepare != nil {
		return errPrepare
	}

	res, errExec := stmt.ExecContext(ctx, user.Id, user.Username, user.Password, user.UpdatedAt, user.CreatedAt)

	if errExec != nil {
		return errExec
	}

	_, errAffected := res.RowsAffected()
	if errAffected != nil {
		return errAffected
	}

	return nil
}

func (repo *AuthRepository) GetUser(ctx context.Context, username string) (entities.User, error) {

	user := entities.User{}

	query := "SELECT id, username, password, updated_at, created_at FROM user WHERE username=?"
	errQuery := repo.db.QueryRowContext(ctx, query, username).Scan(&user.Id, &user.Username, &user.Password, &user.UpdatedAt, &user.CreatedAt)

	if errQuery != nil {
		return entities.User{}, errQuery
	}

	return user, nil
}
