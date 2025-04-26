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

func (s *PostStore) GetPostsRange(start int, length int, admin bool, external bool) (*[]store.Post, error){
    var posts []store.Post
    db := s.db.Preload("Author")

    if !admin{
        if external{
            db = db.Where("external = True")
        }
        db = db.Where("published = True")
    }

    err := db.Order("date desc").Offset(start).Limit(length).Find(&posts).Error

    return &posts, err
}

func (s *PostStore) PatchPost(post store.Post) error {
	updateData := map[string]interface{}{
		"Title":     post.Title,
		"Content":   post.Content,
		"Image":     post.Image,
		"Published": post.Published,
		"External":  post.External,
	}

	return s.db.Model(&store.Post{}).Where("id = ?", post.ID).Updates(updateData).Error
}

func (s *PostStore) DeletePost(postId string) error {
	return s.db.Delete(&store.Post{}, postId).Error
}
