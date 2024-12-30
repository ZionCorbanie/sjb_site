package db

import (
	"errors"
	"fmt"
	"sjb_site/internal/store"
	"os"

	"gorm.io/driver/mysql"	
    "gorm.io/gorm"
)

// open opens a database connection given a database dbName
func open(dbName string) (*gorm.DB, error) {
    dsn := fmt.Sprintf("root:password@tcp(db:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbName)
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
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

	if dbName == "" {
		dbName = "sjb_site"
	}

	db, err := open(dbName)
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&store.User{}, &store.Session{}, &store.Group{}, &store.GroupUser{}, &store.Parent{}, &store.ParentGroup{}, &store.Post{})

	if err != nil {
		panic(err)
	}

	return db
}
