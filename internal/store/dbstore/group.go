package dbstore

import (
	"fmt"
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
	return s.db.Save(group).Error
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

func (s *GroupStore) PatchGroup(group store.Group) error {
	return s.db.Model(&store.Group{}).Where("id = ?", group.ID).Updates(group).Error
}

func (s *GroupStore) DeleteGroup(groupId string) error {
	return s.db.Delete(&store.Group{}, "id = ?", groupId).Error
}

func (s *GroupStore) ValidateInput(name string, groupId uint64) error {
	err := s.db.Model(&store.Group{}).Where("name = ? AND id != ?", name, groupId).First(&store.Group{}).Error
	if err == nil {
		return fmt.Errorf("groep met naam %s bestaat al", name)
	}
	return nil
}
