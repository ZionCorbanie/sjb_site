package store

import "time"

type User struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Email       string    `json:"email" gorm:"type:varchar(255);not null"`
	Password    string    `json:"-" gorm:"type:varchar(255);not null"`
	Username    string    `json:"username" gorm:"type:varchar(255)"`
	FirstName   string    `json:"first_name" gorm:"type:varchar(255)"`
	LastName    string    `json:"last_name" gorm:"type:varchar(255)"`
	StartDate   time.Time `json:"start_date;default:current_timestamp;not null"`
	EndDate     time.Time `json:"end_date"`
	UserType    string    `json:"user_type" gorm:"type:enum('admin','lid','oud_lid');default:lid"`
	Adres       string    `json:"adres" gorm:"type:varchar(255)"`
	PhoneNumber string    `json:"phone_number" gorm:"type:varchar(255)"`
	Image       string    `json:"image" gorm:"type:varchar(255);default:'/static/img/placeholder-150x150.png'"`
}

type Group struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Website     string    `json:"website"`
	GroupType   string    `json:"group_type" gorm:"type:enum('barploeg','bestuur','commissie','gilde','huis','jaarclub','overkoepelend','werkgroep')"`
	StartDate   time.Time `json:"start_date;default:current_timestamp;not null"`
	EndDate     time.Time `json:"end_date"`
	Description string    `json:"description" gorm:"type:varchar(2048)"`
	Image       string    `json:"image" gorm:"type:varchar(255);default:'/static/img/placeholder-group.png'"`
}

type GroupUser struct {
	GroupID  uint   `json:"group_id" gorm:"primaryKey;autoIncrement:false"`
	Group    Group  `gorm:"foreignKey:GroupID;constraint:OnDelete:CASCADE;" json:"group"`
	UserID   uint   `json:"user_id" gorm:"primaryKey;autoIncrement:false"`
	User     User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"user"`
	Status   string `json:"status" gorm:"type:enum('lid','oud_lid','meeloper')"`
	Title    string `json:"title" gorm:"type:varchar(255)"`
	Function string `json:"function" gorm:"type:enum('voorzitter','secretaris','penningmeester')"`
}

type ParentGroup struct {
	ParentID uint  `json:"parent_id"`
	Parent   Group `gorm:"foreignKey:ParentID" json:"parent"`
	ChildID  uint  `json:"child_id"`
	Child    Group `gorm:"foreignKey:ChildID" json:"child"`
}

type Parent struct {
	UserID      uint   `json:"user_id"`
	User        User   `gorm:"foreignKey:UserID" json:"user"`
	Title       string `json:"title" gorm:"type:varchar(255)"`
	Adres       string `json:"adres" gorm:"type:varchar(255)"`
	PhoneNumber string `json:"phone_number" gorm:"type:varchar(255)"`
}

type Post struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `json:"title" gorm:"type:varchar(255)"`
	Content   string    `json:"content"`
	Date      time.Time `json:"date"`
	AuthorID  uint      `json:"author_id"`
	Author    User      `gorm:"foreignKey:AuthorID" json:"author"`
	Published bool      `gorm:"default:False" json:"public"`
	External  bool      `gorm:"default:False" json:"external"`
	Comments  []Comment `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE;" json:"comments"`
}

type Comment struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	Content  string    `json:"content"`
	Date     time.Time `json:"date"`
	AuthorID uint      `json:"author_id"`
	Author   User      `gorm:"foreignKey:AuthorID;constraint:OnDelete:CASCADE;" json:"author"`
	PostID   uint      `json:"post_id"`
}

type Menu struct {
	ID    uint      `gorm:"primaryKey" json:"id"`
	Date  time.Time `json:"date"`
	Name  string    `json:"name" gorm:"type:varchar(255)"`
	Basis string    `json:"basis" gorm:"type:varchar(255)"`
	Vlees string    `json:"vlees" gorm:"type:varchar(255)"`
	Vega  string    `json:"vega" gorm:"type:varchar(255)"`
	Toe   string    `json:"toe" gorm:"type:varchar(255)"`
}

type Session struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	SessionID string `json:"session_id" gorm:"type:varchar(255)"`
	UserID    uint   `json:"user_id"`
	User      User   `gorm:"foreignKey:UserID" json:"user"`
}

type UserStore interface {
	CreateUser(username string, password string) error
	GetUser(username string) (*User, error)
	GetUserById(userId string) (*User, error)
	SearchUsers(search string) ([]*User, error)
	PatchUser(user User) error
	DeleteUser(userId string) error
}

type SessionStore interface {
	CreateSession(session *Session) (*Session, error)
	GetUserFromSession(sessionID string, userID string) (*User, error)
}

type GroupStore interface {
	CreateGroup(group *Group) error
	GetGroupsByType(groupType string) ([]*Group, error)
	GetGroup(groupId string) (*Group, error)
	PatchGroup(group Group) error
	DeleteGroup(groupId string) error
	ValidateInput(name string) error
}

type GroupUserStore interface {
	AddUserToGroup(userId uint, groupId uint) error
	GetUsersByGroup(groupId string) ([]*User, error)
}

type MenuStore interface {
	GetMenu(id string) (*Menu, error)
	GetMenuRange(start string, length string) ([]*Menu, error)
	CreateMenu(menu *Menu) error
}

type PostStore interface {
	CreatePost(post *Post) error
	GetPost(id string) (*Post, error)
	GetPostsRange(start int, length int, admin bool, external bool) ([]*Post, error)
}

type CommentStore interface {
	CreateComment(comment *Comment) error
	GetCommentsByPost(postId string) ([]*Comment, error)
	GetComment(commentId string) (*Comment, error)
	DeleteComment(commentId string) error
}
