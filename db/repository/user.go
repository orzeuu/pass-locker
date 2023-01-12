package repository

import (
	"github.com/orzeuu/pass-locker/password"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Login     string     `json:"login"`
	Password  string     `json:"password"`
	Passwords []Password `json:"passwords" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

type GetUserParams struct {
	Login    string
	Password string
}

type InsertUserParams struct {
	Login    string
	Password string
}

type UserRepository interface {
	GetUser(GetUserParams) (User, error)
	InsertUser(InsertUserParams) (User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db,
	}
}

func (repo *userRepository) GetUser(params GetUserParams) (User, error) {
	var user User
	err := repo.db.First(&user, "login = ?", params.Login).Error
	if err != nil {
		return User{}, err
	}

	if err = password.CheckPassword(params.Password, user.Password); err != nil {
		return User{}, err
	}

	return user, nil
}

func (repo *userRepository) InsertUser(params InsertUserParams) (User, error) {
	hashedPassword, err := password.HashPassword(params.Password)
	if err != nil {
		return User{}, err
	}

	user := User{
		Login:    params.Login,
		Password: hashedPassword,
	}

	err = repo.db.Create(&user).Error

	return user, err
}
