package dbstore

import (
	"sjb_site/internal/store"
	"strconv"

	"gorm.io/gorm"
)


type MenuStore struct {
	db *gorm.DB
}

type NewMenuStoreParams struct {
	DB *gorm.DB
}

func NewMenuStore(params NewMenuStoreParams) *MenuStore {
	return &MenuStore{
		db: params.DB,
	}
}

func (s *MenuStore) GetMenu(menuId string) (*store.Menu, error) {

	var menu store.Menu
	err := s.db.Where("id = ?", menuId).First(&menu).Error

	if err != nil {
		return nil, err
	}
	return &menu, err
}

func (s *MenuStore) GetMenuRange(start string, length string) ([]*store.Menu, error) {
    startInt, err := strconv.ParseInt(start, 10, 32)
    if err != nil {
        return nil, err
    }
    lengthInt, err := strconv.ParseInt(length, 10, 32)
    if err != nil {
        return nil, err
    }

    var menus []*store.Menu
    err = s.db.Where("id >= ?", start).Where("id < ?", startInt+lengthInt).Find(&menus).Error

    return menus, err
}
