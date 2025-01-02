package store

import "time"

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Email        string    `json:"email" gorm:"type:varchar(255);not null"`
	Password     string    `json:"-" gorm:"type:varchar(255);not null"`
	Username     string    `json:"username" gorm:"type:varchar(255)"`
	First_name   string    `json:"first_name" gorm:"type:varchar(255)"`
	Last_name    string    `json:"last_name" gorm:"type:varchar(255)"`
	Start_date   time.Time `json:"start_date"`
	End_date     time.Time `json:"end_date"`
	User_type    string    `json:"user_type" gorm:"type:enum('admin','lid','oud_lid');default:lid"`
	Adres        string    `json:"adres" gorm:"type:varchar(255)"`
	Phone_number string    `json:"phone_number" gorm:"type:varchar(255)"`
	Image        string    `json:"image" gorm:"type:varchar(255);default:'/static/img/placeholder-150x150.png'"`
}

type Parent struct {
	UserID       uint   `json:"user_id"`
	User         User   `gorm:"foreignKey:UserID" json:"user"`
	Title        string `json:"title" gorm:"type:varchar(255)"`
	Adres        string `json:"adres" gorm:"type:varchar(255)"`
	Phone_number string `json:"phone_number" gorm:"type:varchar(255)"`
}

type Group struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `json:"name"`
	Group_type  string    `json:"group_type" gorm:"type:enum('barploeg','bestuur','commissie','gilde','huis','jaarclub','overkoepelend','werkgroep')"`
	Start_date  time.Time `json:"start_date"`
	End_date    time.Time `json:"end_date"`
	Description string    `json:"description" gorm:"type:varchar(2048)"`
	Image       string    `json:"image" gorm:"type:varchar(255);default:'/static/img/placeholder-group.png'"`
}

type GroupUser struct {
	GroupID  uint   `json:"group_id" gorm:"primaryKey;autoIncrement:false"`
	Group    Group  `gorm:"foreignKey:GroupID" json:"group"`
	UserID   uint   `json:"user_id" gorm:"primaryKey;autoIncrement:false"`
	User     User   `gorm:"foreignKey:UserID" json:"user"`
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

type Post struct {
	ID       uint      `gorm:"primaryKey" json:"id"`
	Title    string    `json:"title" gorm:"type:varchar(255)"`
	Content  string    `json:"content"`
	Date     time.Time `json:"date"`
	AuthorID uint      `json:"author_id"`
	Author   User      `gorm:"foreignKey:AuthorID" json:"author"`
}

type Session struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	SessionID string `json:"session_id" gorm:"type:varchar(255)"`
	UserID    uint   `json:"user_id"`
	User      User   `gorm:"foreignKey:UserID" json:"user"`
}

type UserStore interface {
	CreateUser(email string, password string) error
	GetUser(username string) (*User, error)
	GetUserById(userId string) (*User, error)
	SearchUsers(search string) ([]*User, error)
	PatchUser(user User) error
}

type SessionStore interface {
	CreateSession(session *Session) (*Session, error)
	GetUserFromSession(sessionID string, userID string) (*User, error)
}

type GroupStore interface {
	CreateGroup(group *Group) error
	GetGroupsByType(groupType string) ([]*Group, error)
	GetGroup(groupId string) (*Group, error)
}

type GroupUserStore interface {
	AddUserToGroup(userId uint, groupId uint) error
	GetUsersByGroup(groupId string) ([]*User, error)
}
