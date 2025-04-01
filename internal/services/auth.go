package services

import (
	"github.com/evgenyshipko/go-loyality-score-system/internal/logger"
	"github.com/evgenyshipko/go-loyality-score-system/internal/storage"
	"github.com/evgenyshipko/go-loyality-score-system/internal/tokens"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	storage *storage.SQLStorage
}

func NewAuthService(storage *storage.SQLStorage) *AuthService {
	return &AuthService{storage: storage}
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func checkPassword(hashedPwd, inputPwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(inputPwd))
	return err == nil
}

func (a *AuthService) Register(login string, password string) error {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		logger.Instance.Warnf("Ошибка хэширования: %s", err)
		return err
	}

	err = a.storage.InsertUser(login, hashedPassword)
	if err != nil {
		logger.Instance.Warnw("Ошибка инсерта юзера", "err", err.Error())
		return err
	}
	return nil
}

func (a *AuthService) Login(login string, password string) (success bool, err error) {

	user, err := a.storage.GetUser(login)
	if err != nil {
		return false, err
	}

	success = checkPassword(user.HashedPassword, password)

	return success, nil
}

func (a *AuthService) GenerateTokensAndSave(login string) (access string, refresh string, err error) {

	user, err := a.storage.GetUser(login)
	if err != nil {
		return "", "", err
	}

	userId := user.Id

	access, err = tokens.GenerateAccessToken(userId)
	if err != nil {
		return "", "", err
	}
	refresh, err = tokens.GenerateRefreshToken(userId)
	if err != nil {
		return "", "", err
	}
	err = a.storage.SaveUserTokens(userId, access, refresh)
	if err != nil {
		return "", "", err
	}

	return access, refresh, nil
}
