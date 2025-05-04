package dbstore

import (
	"sjb_site/internal/store"
	"time"

	"gorm.io/gorm"
)

type PromoStore struct {
	db           *gorm.DB
}

type NewPromoStoreParams struct {
	DB           *gorm.DB
}

func NewPromoStore(params NewPromoStoreParams) *PromoStore {
	return &PromoStore{
		db: params.DB,
	}
}

func (s *PromoStore) SavePromo(promo *store.Promo) error {
	return s.db.Save(promo).Error
}

func (s *PromoStore) GetActivePromos() (*[]store.Promo, error) {
	var promos []store.Promo
	err := s.db.Where("start_date <= ? AND end_date >= ?", time.Now(), time.Now()).Find(&promos).Error

	if err != nil {
		return nil, err
	}
	return &promos, err
}

func (s *PromoStore) GetAllPromos() (*[]store.Promo, error) {
	var promos []store.Promo
	err := s.db.Find(&promos).Error

	if err != nil {
		return nil, err
	}
	return &promos, err
}

func (s *PromoStore) GetPromo(promoId string) (*store.Promo, error) {
	var promo store.Promo
	err := s.db.Where("id = ?", promoId).First(&promo).Order("start_date DESC").Error

	if err != nil {
		return nil, err
	}
	return &promo, err
}

func (s *PromoStore) PatchPromo(promo store.Promo) error {
	return s.db.Model(&store.Promo{}).Where("id = ?", promo.ID).Updates(promo).Error
}

func (s *PromoStore) DeletePromo(promoId string) error {
	return s.db.Delete(&store.Promo{}, promoId).Error
}

func (s *PromoStore) DeleteInactivePromos() error {
	return s.db.Where("end_date < ?", time.Now()).Delete(&store.Promo{}).Error
}
