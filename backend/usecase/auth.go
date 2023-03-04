package usecase

import (
	"errors"
	"time"

	. "github.com/arumakan1727/todo-app-go-react/domain"
	"github.com/google/uuid"
)

type authUc struct {
	authTokenMaxAge time.Duration
	kvs             KVS
	userRepo        UserReadRepository
}

type UserReadRepository interface {
	GetUserByEmail(Ctx, string) (User, error)
}

func NewAuthUsecase(
	userRepo UserReadRepository,
	kvs KVS,
	authTokenMaxAge time.Duration,
) AuthUsecase {
	return &authUc{
		authTokenMaxAge: authTokenMaxAge,
		kvs:             kvs,
		userRepo:        userRepo,
	}
}

func (uc *authUc) GetAuthTokenMaxAge() time.Duration {
	return uc.authTokenMaxAge
}

func (uc *authUc) IssueAuthToken(ctx Ctx, email, passwd string) (AuthToken, error) {
	user, err := uc.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return "", ErrIncorrectEmailOrPasswd
		}
		return "", err
	}

	if err := ComparePassword(user.PasswdHash, email, passwd); err != nil {
		return "", ErrIncorrectEmailOrPasswd
	}

	ruuid, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	token := AuthToken(ruuid.String())
	if err := uc.kvs.SaveAuth(ctx, token, user.ID, uc.authTokenMaxAge); err != nil {
		return "", err
	}
	return token, nil
}

func (uc *authUc) ValidateAuthToken(ctx Ctx, token AuthToken) (UserID, error) {
	uid, err := uc.kvs.FetchAuth(ctx, token)
	if err != nil {
		return 0, ErrUnauthorized
	}
	return uid, nil
}

func (uc *authUc) RevokeAuthToken(ctx Ctx, token AuthToken) error {
	err := uc.kvs.DeleteAuth(ctx, token)
	return err
}
