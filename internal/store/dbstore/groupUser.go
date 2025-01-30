package dbstore

import (
	"sjb_site/internal/store"
	"time"

	"gorm.io/gorm"
)

type GroupUserStore struct {
	db *gorm.DB
}

type NewGroupUserStoreParams struct {
	DB *gorm.DB
}

func NewGroupUserStore(params NewGroupUserStoreParams) *GroupUserStore {
	return &GroupUserStore{
		db: params.DB,
	}
}

func (s *GroupUserStore) AddUserToGroup(userId, groupId uint) error {
	return s.db.Create(&store.GroupUser{
		UserID:    userId,
		GroupID:   groupId,
		Status:    "lid",
		StartDate: time.Now(),
	}).Error
}

func (s *GroupUserStore) GetUsersByGroup(groupId string) ([]*store.User, error) {
	var users []*store.User

	err := s.db.Joins("JOIN group_users ON group_users.user_id = users.id").
		Where("group_users.group_id = ?", groupId).
		Find(&users).Error

	return users, err
}

func (s *GroupUserStore) GetGroupUsersByGroup(groupId string) ([]*store.GroupUser, error) {
	var members []*store.GroupUser

	err := s.db.Preload("User").Where("group_id = ?", groupId).Find(&members).Error
	return members, err
}

func (s *GroupUserStore) DeleteGroupUser(userId, groupId uint) error {
	return s.db.Delete(&store.GroupUser{}, "user_id = ? AND group_id = ?", userId, groupId).Error
}

func (s *GroupUserStore) UpdateGroupUser(groupUser store.GroupUser) error {
	return s.db.Model(&store.GroupUser{}).Where("user_id = ? AND group_id = ?", groupUser.UserID, groupUser.GroupID).Updates(groupUser).Error
}
