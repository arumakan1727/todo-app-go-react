package usecase

import (
	"fmt"

	. "github.com/arumakan1727/todo-app-go-react/domain"
)

type userUc struct {
	repo Repository
}

func NewUserUsecase(repo Repository) UserUsecase {
	return &userUc{
		repo: repo,
	}
}

func (uc *userUc) Store(ctx Ctx, email, passwd, displayName string) (User, error) {
	pwhash, err := HashPassword(email, passwd)
	if err != nil {
		return User{}, fmt.Errorf("cannot hash password: %w", err)
	}
	user := User{
		Role:        "user",
		Email:       email,
		PasswdHash:  pwhash,
		DisplayName: displayName,
	}
	if err := uc.repo.StoreUser(ctx, &user); err != nil {
		return User{}, err
	}
	return user, nil
}

func (uc *userUc) List(ctx Ctx) ([]User, error) {
	return uc.repo.ListUsers(ctx)
}
