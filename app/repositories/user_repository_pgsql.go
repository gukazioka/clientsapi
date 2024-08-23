package repositories

import (
	"context"
	"fmt"
	"os"

	"github.com/gkazioka/clientsapi/app/domain"
	"github.com/gkazioka/clientsapi/app/infra/config"
	"github.com/gkazioka/clientsapi/app/infra/database"
	"github.com/gkazioka/clientsapi/app/types"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepositoryPostgres struct {
	dbPool *pgxpool.Pool
}

func (u *UserRepositoryPostgres) UserExists(ctx context.Context, user domain.User) (bool, error) {
	var found bool
	err := u.dbPool.QueryRow(ctx, database.FindUserQuery, user.Code).Scan(&found)

	if err != nil {
		return false, err
	}
	return found, nil
}

func (u *UserRepositoryPostgres) FindUserByCode(ctx context.Context, userCode string) (*domain.User, error) {
	var userFound domain.User
	err := u.dbPool.QueryRow(ctx, database.FindUserByCodeQuery, userCode).Scan(&userFound.Name, &userFound.Code)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &userFound, nil
}

func (u *UserRepositoryPostgres) Save(ctx context.Context, user domain.User) error {
	exists, error := u.UserExists(ctx, user)

	if error != nil {
		return fmt.Errorf("internal error")
	}

	if exists {
		return types.ErrorAlreadyExists
	}

	_, err := u.dbPool.Query(ctx, database.SaveUserQuery, user.Name, user.Code)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Write error: %v\n", err)
		return err
	}
	return nil
}

func (u *UserRepositoryPostgres) ListAll(ctx context.Context) []domain.User {
	var users []domain.User
	var user domain.User

	users = make([]domain.User, 0)

	rows, err := u.dbPool.Query(ctx, database.ListUsersQuery)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Read error: %v\n", err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&user.Name, &user.Code); err != nil {
			fmt.Fprintf(os.Stderr, "Read error: %v\n", err)
		}
		users = append(users, user)
	}

	return users
}

func NewUserRepositoryPostgres() *UserRepositoryPostgres {
	return &UserRepositoryPostgres{dbPool: database.NewDatabase(*config.GetConfig())}
}
