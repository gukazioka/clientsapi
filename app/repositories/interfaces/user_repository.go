package interfaces

import (
	"context"

	"github.com/gkazioka/clientsapi/app/domain"
)

type UserRepository interface {
	ListAll(ctx context.Context) []domain.User
	Save(ctx context.Context, user domain.User) error
	FindUserByCode(ctx context.Context, userCode string) (*domain.User, error)
}
