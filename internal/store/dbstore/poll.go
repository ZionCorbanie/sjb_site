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
    err := s.db.First(&poll, pollId).Error

    return &poll, err
}

func (s *PollStore) DeletePoll(pollId string) error {
    return s.db.Delete(&store.Poll{}, pollId).Error
}

func (s *PollStore) PatchPoll(poll store.Poll) error {
    updateData := map[string]interface{}{
        "Title": poll.Title,
        "Options": poll.Options,
    }

    return s.db.Model(&poll).Updates(updateData).Error
}
