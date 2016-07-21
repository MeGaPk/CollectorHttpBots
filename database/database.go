package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Bot struct {
	gorm.Model
	Link string
	Header string
	Body string
	Form string
	PostForm string
}


type DatabaseConnection struct {
	*gorm.DB
}
//TODO: https://github.com/mattn/go-sqlite3/issues/214
func New(databasePath string) *DatabaseConnection {
	db, err := gorm.Open("sqlite3", databasePath)
	if err != nil {
		panic("failed to connect database")
	}
	// Migrate the schema
	db.AutoMigrate(&Bot{})
	return &DatabaseConnection{db}
}

func (db *DatabaseConnection) GetUrls() (bots []Bot){
	db.Find(&bots)
	return bots;
}

func (db *DatabaseConnection) AddUrl(bot *Bot) {
	db.Create(bot)
}
