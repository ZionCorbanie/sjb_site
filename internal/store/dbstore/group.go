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

func (s *GroupStore) GetSimelarGroups(group *store.Group) (*[]store.Group, string, error) {
    var groups []store.Group
    switch group.GroupType {   
    case "jaarclub":
        err := s.db.Where("group_type = ? AND year(start_date) = year(?) AND id != ?", "jaarclub", group.StartDate, group.ID).Find(&groups).Error
        if err != nil {
            return nil, "", err
        }
        return &groups, "Jaarclubs uit " + group.StartDate.Format("2006"), err
    case "barploeg":
        err := s.db.Where("group_type = ? AND end_date IS NULL AND id != ?", "barploeg", group.ID).Find(&groups).Error
        if err != nil {
            return nil, "", err
        }
        return &groups, "Barbloegen", err
    case "commissie":
        err := s.db.Joins("JOIN parent_groups pg ON groups.id = pg.child_id JOIN groups p on p.id = pg.parent_id").Where("groups.group_type = ? AND groups.id != ?", "commissie", group.ID).Find(&groups).Error
        if err != nil {
            return nil, "", err
        }
        return &groups, "Commissies", err
    case "bestuur":
        err := s.db.Where("group_type = ? AND id != ?", "bestuur", group.ID).Find(&groups).Error
        if err != nil {
            return nil, "", err
        }
        return &groups, "Besturen", err
    case "huis":
        err := s.db.Where("group_type = ? AND id != ?", "huis", group.ID).Find(&groups).Error
        if err != nil {
            return nil, "", err
        }
        return &groups, "Huizen", err
    case "gilde":
        err := s.db.Where("group_type = ? AND end_date IS NULL AND id != ?", "gilde", group.ID).Find(&groups).Error
        if err != nil {
            return nil, "", err
        }
        return &groups, "Gilden", err
    case "werkgroep":
        err := s.db.Where("group_type = ? AND end_date IS NULL AND id != ?", "werkgroep", group.ID).Find(&groups).Error
        if err != nil {
            return nil, "", err
        }
    }

    return nil, "", nil
}







