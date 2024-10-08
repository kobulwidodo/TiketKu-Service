package user

import (
	"errors"
	userDom "go-clean/src/business/domain/user"
	"go-clean/src/business/entity"
	"go-clean/src/lib/auth"

	"golang.org/x/crypto/bcrypt"
)

type Interface interface {
	Create(params entity.CreateUserParam) (entity.User, error)
	Login(params entity.LoginUserParam) (string, error)
	GetById(id uint) (entity.User, error)
}

type user struct {
	user userDom.Interface
	auth auth.Interface
}

func Init(ad userDom.Interface, auth auth.Interface) Interface {
	a := &user{
		user: ad,
		auth: auth,
	}

	return a
}

func (a *user) Create(params entity.CreateUserParam) (entity.User, error) {
	user := entity.User{
		Email: params.Email,
		Name:  params.Name,
	}

	hashPass, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.MinCost)
	if err != nil {
		return user, err
	}

	user.Password = string(hashPass)

	newUser, err := a.user.Create(user)
	if err != nil {
		return newUser, err
	}

	return newUser, nil
}

func (a *user) GetById(id uint) (entity.User, error) {
	user, err := a.user.Get(entity.UserParam{
		ID: id,
	})
	if err != nil {
		return user, err
	}

	return user, nil
}

func (a *user) Login(params entity.LoginUserParam) (string, error) {
	user, err := a.user.Get(entity.UserParam{
		Email: params.Email,
	})
	if err != nil {
		return "", err
	}

	if user.ID == 0 {
		return "", errors.New("user tidak ditemukan atau password tidak sesuai")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password)); err != nil {
		return "", errors.New("user tidak ditemukan atau password tidak sesuai")
	}

	token, err := a.auth.GenerateToken(user.ConvertToAuthUser())
	if err != nil {
		return "", err
	}

	return token, nil
}
