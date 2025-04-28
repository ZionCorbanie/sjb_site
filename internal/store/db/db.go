package db

import (
	"errors"
	"fmt"
	"os"
	"sjb_site/internal/store"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// open opens a database connection given a database dbName
func open(dbName string) (*gorm.DB, error) {

	password :=  os.Getenv("DB_PASSWORD")
	user := os.Getenv("DB_USER")
	server := os.Getenv("DB")
	//port := os.Getenv("DB_PORT")

	if dbName == "" {
		dbName = "sjb_site"
	}

	
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, server, dbName)

	config := &gorm.Config{}

	if os.Getenv("env") == "production" {
		config.Logger = logger.Default.LogMode(logger.Silent)
	}

	db, err := gorm.Open(mysql.Open(dsn), config)
	if err != nil {
		return nil, errors.Join(err, errors.New("failed to open database"))
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	err = sqlDB.Ping()
	if err != nil {
		return nil, err
	}

	// make the temp directory if it doesn't exist
	err = os.MkdirAll("/tmp", 0755)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// MustOpen opens a database connection and panics if it fails
func MustOpen(dbName string) *gorm.DB {

	db, err := open(dbName)
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&store.User{}, &store.Session{}, &store.Group{}, &store.GroupUser{}, &store.Parent{}, &store.ParentGroup{}, &store.Post{}, &store.Menu{}, &store.Comment{}, &store.Poll{}, &store.PollOption{}, &store.PollVote{}, &store.CalendarItem{})

	if err != nil {
		panic(err)
	}

	return db
}
