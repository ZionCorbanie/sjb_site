package dbstore

import (
	"sjb_site/internal/store"
	"time"

	"gorm.io/gorm"
)

type CalendarStore struct {
    db *gorm.DB
}

type NewCalendarStoreParams struct {
    DB *gorm.DB
}

func NewCalendarStore(params NewCalendarStoreParams) *CalendarStore {
    return &CalendarStore{
        db: params.DB,
    }
}

func (s *CalendarStore) CreateCalendarItem(event *store.CalendarItem) error {
    return s.db.Create(event).Error
}

func (s *CalendarStore) GetCalendarItems(day int) (*[]store.CalendarItem, error) {
	var targetDate time.Time

	err := s.db.
		Model(&store.CalendarItem{}).
		Select("start_date").
		Group("DATE_FORMAT(start_date, '%Y-%m-01')").
		Order("start_date ASC").
		Limit(1).
		Offset(day).
		Pluck("start_date", &targetDate).Error
	if err != nil {
		return nil, err
	}

	baseDate := targetDate.Truncate(24 * time.Hour)
	endOfDay := baseDate.Add(24 * time.Hour)

	var items []store.CalendarItem

	err = s.db.
		Where("start_date > ? AND start_date <?", baseDate, endOfDay).
		Order("start_date ASC").
		Find(&items).Error
	if err != nil {
		return nil, err
	}

	return &items, nil
}

func (s *CalendarStore) GetCalendarItem(id string) (*store.CalendarItem, error) {
	var item store.CalendarItem
	err := s.db.Where("id = ?", id).First(&item).Error

	return &item, err
}
