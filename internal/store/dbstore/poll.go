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
    err := s.db.Order("id DESC").Find(&polls).Error
    return polls, err
}

func (s *PollStore) DeletePoll(pollId string) error {
    return s.db.Delete(&store.Poll{}, pollId).Error
}

func (s *PollStore) Activate(pollId string) error {
    return s.db.Exec("UPDATE polls SET active = (CASE WHEN id = ? THEN 1 ELSE 0 END)", pollId).Error
}

func (s *PollStore) GetActivePoll() (*store.Poll, error) {
    var poll store.Poll
    err := s.db.Where("active = 1").First(&poll).Error
    return &poll, err
}

func (s *PollStore) PutPoll(poll store.Poll) error {
    tx := s.db.Begin()

    if err := tx.Where("poll_id = ?", poll.ID).Delete(&store.PollOption{}).Error; err != nil {
        tx.Rollback()
        return err
    }

    if err := tx.Model(&store.Poll{}).Where("id = ?", poll.ID).
        Updates(map[string]interface{}{
            "title": poll.Title, // Add other fields if necessary
        }).Error; err != nil {
        tx.Rollback()
        return err
    }

    for i := range poll.Options {
        poll.Options[i].PollID = poll.ID // Ensure PollID is set
    }

    if len(poll.Options) > 0 {
        if err := tx.Create(&poll.Options).Error; err != nil {
            tx.Rollback()
            return err
        }
    }

    // Commit transaction
    return tx.Commit().Error
}

func (s *PollStore) Vote(pollId uint, optionId uint, userId uint) error {
    tx := s.db.Begin()

    if err := tx.Where("option_id = ? AND user_id = ?", optionId, userId).Delete(&store.PollVote{}).Error; err != nil {
        tx.Rollback()
        return err
    }

    if err := tx.Create(&store.PollVote{
        PollID:  pollId,
        OptionID:   optionId,
        UserID:   userId,
        
    }).Error; err != nil {
        tx.Rollback()
        return err
    }

    return tx.Commit().Error
}

func (s *PollStore) GetPollVotes(pollID uint, userID uint) (*store.Poll, bool) {
    var poll store.Poll

    s.db.Preload("Options").First(&poll, pollID)

    var voteCounts []struct {
        OptionID  uint
        VoteCount int
    }

    s.db.Table("poll_votes").
        Select("poll_votes.option_id, COUNT(poll_votes.id) as vote_count").
        Where("poll_votes.poll_id = ?", pollID).
        Group("poll_votes.option_id").
        Scan(&voteCounts)

    voteCountMap := make(map[uint]int)
    for _, v := range voteCounts {
        voteCountMap[v.OptionID] = v.VoteCount
    }
    for i := range poll.Options {
        poll.Options[i].VoteCount = voteCountMap[poll.Options[i].ID]
    }

    var userVoteExists bool
    s.db.Table("poll_votes").
        Where("poll_id = ? AND user_id = ?", pollID, userID).
        Select("1").
        Limit(1).
        Scan(&userVoteExists)

    return &poll, userVoteExists
}

func (s *PollStore) DeleteVote(pollId uint, userId uint) error{
    return s.db.Where("poll_id = ? AND user_id = ?", pollId, userId).Delete(&store.PollVote{}).Error
}
