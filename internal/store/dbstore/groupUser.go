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

	err := s.db.Joins("JOIN group_users ON group_users.user_id = users.id").
		Where("group_users.group_id = ?", groupId).
		Find(&users).Error

	return users, err
}

func (s *GroupUserStore) GetGroupsByUser(userId string) ([]*store.Group, error) {
    var groups []*store.Group

    err := s.db.Joins("JOIN group_users ON group_users.group_id = groups.id").
        Where("group_users.user_id = ?", userId).
        Order("groups.group_type").
        Find(&groups).Error

    return groups, err
}
