package repository

import (
	"gorm.io/gorm"
)

type Password struct {
	gorm.Model
	Item     string `json:"item" gorm:"not null"`
	Login    string `json:"login" gorm:"not null,unique"`
	Password string `json:"password" gorm:"not null"`
	UserId   uint   `json:"userId" gorm:"not null"`
}

type AddPasswordParams struct {
	Item     string
	Login    string
	Password string
	UserId   uint
}

type GetPasswordParams struct {
	Item string
}

type GetAllPasswordsParams struct {
	UserId uint
}

type DeletePasswordParams struct {
	ID uint64
}

type PasswordRepository interface {
	AddPassword(AddPasswordParams) (Password, error)
	GetPassword(GetPasswordParams) (Password, error)
	GetAllPasswords(GetAllPasswordsParams) ([]Password, error)
	DeletePassword(DeletePasswordParams)
}

type passwordRepository struct {
	db *gorm.DB
}

func NewPasswordRepository(db *gorm.DB) PasswordRepository {
	return &passwordRepository{
		db,
	}
}

func (repo *passwordRepository) AddPassword(params AddPasswordParams) (Password, error) {
	password := Password{
		Item:     params.Item,
		Login:    params.Login,
		Password: params.Password,
		UserId:   params.UserId,
	}
	err := repo.db.Create(&password).Error
	return password, err
}

func (repo *passwordRepository) GetPassword(params GetPasswordParams) (Password, error) {
	var password Password

	err := repo.db.First(&password, "item = ?", params.Item).Error

	return password, err
}

func (repo *passwordRepository) GetAllPasswords(params GetAllPasswordsParams) ([]Password, error) {
	var passwords []Password
	err := repo.db.Where("user_id = ?", params.UserId).Find(&passwords).Error

	return passwords, err
}

func (repo *passwordRepository) DeletePassword(params DeletePasswordParams) {
	repo.db.Delete(&Password{}, params.ID)
}
