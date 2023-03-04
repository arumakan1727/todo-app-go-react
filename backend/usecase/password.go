package usecase

import (
	"strings"
	"unsafe"

	"github.com/arumakan1727/todo-app-go-react/domain"
	"golang.org/x/crypto/bcrypt"
)

func makePasswdHashPayload(email, passwd string) []byte {
	const salt = "iVe$o5uGhe,x1yeetoo^P9ohPhoh3AhbaeHohde9"
	s := passwd + strings.ToLower(email) + salt
	b := unsafe.Slice(unsafe.StringData(s), len(s))

	// bcrypt.GenerateFromPassword()が72byteまでしか受けつけない
	return b[:72]
}

func HashPassword(email, passwd string) ([]byte, error) {
	p := makePasswdHashPayload(email, passwd)
	return bcrypt.GenerateFromPassword(p, bcrypt.DefaultCost)
}

func ComparePassword(hashedPasswd []byte, email, passwd string) error {
	p := makePasswdHashPayload(email, passwd)
	err := bcrypt.CompareHashAndPassword(hashedPasswd, p)
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return domain.ErrIncorrectEmailOrPasswd
		}
		return err
	}
	return nil
}
