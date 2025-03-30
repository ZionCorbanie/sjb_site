package dbstore

import (
	"sjb_site/internal/store"

	"gorm.io/gorm"
)

type PollStore struct {
    db *gorm.DB
}

type NewPollStoreParams struct {
    DB *gorm.DB
}

func NewPollStore(params NewPollStoreParams) *PollStore {
    return &PollStore{
        db: params.DB,
    }
}

func (s *PollStore) CreatePoll(poll *store.Poll) error {
    return s.db.Save(poll).Error
}

func (s *PollStore) GetPoll(pollId string) (*store.Poll, error) {
    var poll store.Poll
    err := s.db.Preload("Options").First(&poll, pollId).Error

    return &poll, err
}

func (s *PollStore) GetPolls() ([]*store.Poll, error) {
    var polls []*store.Poll
    err := s.db.Find(&polls).Error
    return polls, err
}

func (s *PollStore) DeletePoll(pollId string) error {
    return s.db.Delete(&store.Poll{}, pollId).Error
}

func (s *PollStore) PutPoll(poll store.Poll) error {
    //remove assodiated options
    err := s.db.Where("poll_id = ?", poll.ID).Delete(store.PollOption{}).Error
    err = s.db.Save(&poll).Error
    return err
}
