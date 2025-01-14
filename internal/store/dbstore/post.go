package dbstore

import (
	"sjb_site/internal/store"

	"gorm.io/gorm"
)


type PostStore struct {
	db *gorm.DB
}

type NewPostStoreParams struct {
	DB *gorm.DB
}

func NewPostStore(params NewPostStoreParams) *PostStore {
	return &PostStore{
		db: params.DB,
	}
}

func (s *PostStore) CreatePost(post *store.Post) error{
    return s.db.Save(post).Error
}

func (s *PostStore) GetPost(postId string) (*store.Post, error){
	var post store.Post
	err := s.db.Preload("Author").Where("id = ?", postId).First(&post).Error

	if err != nil {
		return nil, err
	}
	return &post, err
}

func (s *PostStore) GetPostsRange(start int, length int) ([]*store.Post, error){
    var posts []*store.Post
    err := s.db.Preload("Author").Order("date desc").Offset(start).Limit(length).Find(&posts).Error

    return posts, err
}
