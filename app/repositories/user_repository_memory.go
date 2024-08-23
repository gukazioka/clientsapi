package repositories

import (
	"context"

	"github.com/gkazioka/clientsapi/app/domain"
)

type UserRepositoryMemory struct {
	users []domain.User
}

func (u *UserRepositoryMemory) Save(ctx context.Context, user domain.User) error {
	u.users = append(u.users, user)
	return nil
}

func (u UserRepositoryMemory) ListAll(ctx context.Context) []domain.User {
	return u.users
}

func (u UserRepositoryMemory) FindUserByCode(ctx context.Context, userCode string) (*domain.User, error) {
	for _, user := range u.users {
		if user.Code == userCode {
			return &user, nil
		}
	}
	return nil, nil
}
