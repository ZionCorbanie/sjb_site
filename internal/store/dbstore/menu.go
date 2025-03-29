package dbstore

import (
	"sjb_site/internal/store"
	"strconv"
	"time"

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

func (s *MenuStore) GetMenu(menuId string) (*store.Menu) {
	var menu store.Menu
    err := s.db.Where("id = ?", menuId).First(&menu).Error

	if err != nil {
        menu = store.Menu{}
        menuIdInt, _ := strconv.ParseInt(menuId, 10, 64)
        menu.ID = uint(menuIdInt)
        menu.Name = "Onbekend"

        menu.Date = time.Unix(menuIdInt*60*60*24, 0)
        if menu.Date.Weekday() == time.Saturday || menu.Date.Weekday() == time.Sunday {
            menu.Name = "Eettafel gesloten"
        }
    }

	return &menu
}

func (s *MenuStore) GetMenuRange(start int, length int) ([]*store.Menu, error) {
    menus := make([]*store.Menu, length)

    for i := 0; i < length; i++ {
        var menu store.Menu
        _ = s.db.Where("id = ?", start+i).Find(&menu).Error
        menus[i] = &menu
    }

    return menus, nil
}

func (s *MenuStore) CreateMenu(menu *store.Menu) error {
    return s.db.Save(menu).Error
}
