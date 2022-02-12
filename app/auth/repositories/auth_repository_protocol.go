package repositories

import (
	"TodoAPI/app/auth/entities"
	"context"
)

type AuthRepositoryProtocol interface {
	Store(ctx context.Context, user *entities.User) error
	GetUser(ctx context.Context, username string) (entities.User, error)
}
