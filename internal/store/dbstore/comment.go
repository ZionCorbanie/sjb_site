package dbstore

import (
	"sjb_site/internal/store"

	"gorm.io/gorm"
)

type CommentStore struct {
    db *gorm.DB
}

type NewCommentStoreParams struct {
    DB *gorm.DB
}

func NewCommentStore(params NewCommentStoreParams) *CommentStore {
    return &CommentStore{
        db: params.DB,
    }
}

func (s *CommentStore) CreateComment(comment *store.Comment) error {
    return s.db.Create(comment).Error
}

func (s *CommentStore) GetComment(commentId string) (*store.Comment, error) {
    var comment store.Comment
    err := s.db.Preload("Author").First(&comment, commentId).Error

    return &comment, err
}

func (s *CommentStore) GetCommentsByPost(postId string) (*[]store.Comment, error) {
    var comments []store.Comment
    err := s.db.Preload("Author").Where("post_id = ?", postId).Find(&comments).Error

    return &comments, err
}

func (s *CommentStore) DeleteComment(commentId string) error {
    return s.db.Delete(&store.Comment{}, commentId).Error
}
