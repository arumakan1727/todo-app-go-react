package usecase

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"unsafe"

	"github.com/arumakan1727/todo-app-go-react/domain"
	. "github.com/arumakan1727/todo-app-go-react/domain"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userUc struct {
	authTokenLifeSpan time.Duration
	kvs               KVS
	repo              Repository
}

func NewUserUsecase(repo Repository, kvs KVS, authTokenLifeSpan time.Duration) UserUcase {
	return &userUc{
		authTokenLifeSpan: authTokenLifeSpan,
		kvs:               kvs,
		repo:              repo,
	}
}

func (uc *userUc) makePasswdHashPayload(email, passwd string) []byte {
	const salt = "iVe$o5uGhe,x1yeetoo^P9ohPhoh3AhbaeHohde9"
	s := passwd + strings.ToLower(email) + salt
	b := unsafe.Slice(unsafe.StringData(s), len(s))

	// bcrypt.GenerateFromPassword()が72byteまでしか受けつけない
	return b[:72]
}

func (uc *userUc) hashPassword(email, passwd string) ([]byte, error) {
	p := uc.makePasswdHashPayload(email, passwd)
	return bcrypt.GenerateFromPassword(p, bcrypt.DefaultCost)
}

func (uc *userUc) comparePassword(hashedPasswd []byte, email, passwd string) error {
	p := uc.makePasswdHashPayload(email, passwd)
	err := bcrypt.CompareHashAndPassword(hashedPasswd, p)
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return domain.ErrIncorrectEmailOrPasswd
		}
		return err
	}
	return nil
}

func (uc *userUc) Store(ctx Ctx, email, passwd, displayName string) (User, error) {
	pwhash, err := uc.hashPassword(email, passwd)
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

func (uc *userUc) GetAuthTokenLifeSpan() time.Duration {
	return uc.authTokenLifeSpan
}

func (uc *userUc) IssueAuthToken(ctx Ctx, email, passwd string) (AuthToken, error) {
	user, err := uc.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return "", domain.ErrIncorrectEmailOrPasswd
		}
		return "", err
	}

	if err := uc.comparePassword(user.PasswdHash, email, passwd); err != nil {
		return "", domain.ErrIncorrectEmailOrPasswd
	}

	ruuid, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	token := AuthToken(ruuid.String())
	if err := uc.kvs.SaveAuth(ctx, token, user.ID, uc.authTokenLifeSpan); err != nil {
		return "", err
	}
	return token, nil
}

func (uc *userUc) ValidateAuthToken(ctx Ctx, token AuthToken) (UserID, error) {
	uid, err := uc.kvs.FetchAuth(ctx, token)
	if err != nil {
		return 0, domain.ErrUnauthorized
	}
	return uid, nil
}

func (uc *userUc) RevokeAuthToken(ctx Ctx, token AuthToken) error {
	err := uc.kvs.DeleteAuth(ctx, token)
	return err
}
