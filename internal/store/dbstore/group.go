package dbstore

import (
	"sjb_site/internal/store"

	"gorm.io/gorm"
)

type GroupStore struct {
	db *gorm.DB
}

type NewGroupStoreParams struct {
	DB *gorm.DB
}

func NewGroupStore(params NewGroupStoreParams) *GroupStore {
	return &GroupStore{
		db: params.DB,
	}
}

func (s *GroupStore) CreateGroup(group *store.Group) error {
	return s.db.Create(group).Error
}

func (s *GroupStore) GetGroup(groupId string) (*store.Group, error) {

	var group store.Group
	err := s.db.Where("id = ?", groupId).First(&group).Error

	if err != nil {
		return nil, err
	}
	return &group, err
}

func (s *GroupStore) GetGroupsByType(groupType string) ([]*store.Group, error) {
	var groups []*store.Group
	err := s.db.Where("group_type = ?", groupType).Find(&groups).Error

	if err != nil {
		return nil, err
	}
	return groups, err
}

func (s *GroupStore) GetJaarclubs(group *store.Group) (*[]store.Group, error) {
    var groups []store.Group
    err := s.db.Where("group_type = ? AND year(start_date) = year(?) AND id != ?", "jaarclub", group.StartDate, group.ID).Find(&groups).Error

    if err != nil {
        return nil, err
    }
    return &groups, err
}
