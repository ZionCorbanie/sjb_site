package dbstore

import (
	"fmt"
	"sjb_site/internal/hash"
	"sjb_site/internal/store"

	"gorm.io/gorm"
)

type UserStore struct {
	db           *gorm.DB
	passwordhash hash.PasswordHash
}

type NewUserStoreParams struct {
	DB           *gorm.DB
	PasswordHash hash.PasswordHash
}

func NewUserStore(params NewUserStoreParams) *UserStore {
	return &UserStore{
		db:           params.DB,
		passwordhash: params.PasswordHash,
	}
}

func (s *UserStore) CreateUser(username string, password string) error {

	hashedPassword, err := s.passwordhash.GenerateFromPassword(password)
	if err != nil {
		return err
	}

	return s.db.Create(&store.User{
		Username: username,
		Password: hashedPassword,
	}).Error
}

func (s *UserStore) GetUser(username string) (*store.User, error) {

	var user store.User
	err := s.db.Where("username = ?", username).First(&user).Error

	if err != nil {
		return nil, err
	}
	return &user, err
}

func (s *UserStore) GetUserById(userId string) (*store.User, error) {

	var user store.User
	err := s.db.Where("id = ?", userId).First(&user).Error

	if err != nil {
		return nil, err
	}
	return &user, err
}

func (s *UserStore) SearchUsers(search string) ([]*store.User, error) {

	var users []*store.User
	err := s.db.Where("username like ?", fmt.Sprintf("%%%s%%", search)).Find(&users).Error

	if err != nil {
		return nil, err
	}
	return users, err
}

func (s *UserStore) PatchUser(user store.User) error {
	return s.db.Model(&store.User{}).Where("id = ?", user.ID).Updates(user).Error
}

func (s *UserStore) DeleteUser(userId string) error {
	return s.db.Delete(&store.User{}, userId).Error
}
