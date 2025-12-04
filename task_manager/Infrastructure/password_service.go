package Infrastructure

import "golang.org/x/crypto/bcrypt"

type PasswordService interface {
	Hash(password string) (string, error)
	Compare(hash, password string) error
}

type BcryptPasswordService struct{}

func NewBcryptPasswordService() PasswordService {
	return &BcryptPasswordService{}
}

func (b *BcryptPasswordService) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (b *BcryptPasswordService) Compare(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
