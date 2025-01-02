package dbstore

import (
	"sjb_site/internal/store"

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

func (s *GroupUserStore) AddUserToGroup(userId uint, groupId uint) error {
	return s.db.Create(&store.GroupUser{
		UserID:  userId,
		GroupID: groupId,
	}).Error
}

func (s *GroupUserStore) GetUsersByGroup(groupId string) ([]*store.User, error) {
	var users []*store.User

	// err := s.db.Table("group_users").Select("user_id").Where("group_id = ?", groupId).Find(&users).Error
	// err := s.db.Where("group_id=?", groupId).Find(&store.GroupUser{}).Scan(&users).Error
	err := s.db.Joins("JOIN group_users ON group_users.user_id = users.id").
		Where("group_users.group_id = ?", groupId).
		Find(&users).Error

	return users, err
}
