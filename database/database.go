package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"fmt"
)

type Bot struct {
	gorm.Model
	Link     string
	Header   string
	Body     string
	Form     string
	PostForm string
	RemoteIp string
}

type PasteText struct {
	gorm.Model
	Code string `gorm:"primary_key"`
	Text string `gorm:"size:999999999"`
}

type DatabaseConnection struct {
	*gorm.DB
}
//TODO: https://github.com/mattn/go-sqlite3/issues/214
func NewSqlite3(databasePath string) *DatabaseConnection {
	db, err := gorm.Open("sqlite3", databasePath)
	if err != nil {
		panic("failed to connect database sqlite, error: " + err.Error())
	}
	// Migrate the schema
	return &DatabaseConnection{db}
}

func NewMySQL(host string, port int, login string, password string, dbname string) *DatabaseConnection {
	scheme_url := fmt.Sprintf("%v:%v@tcp(%v:%d)/%v?charset=utf8&parseTime=True&loc=Local", login, password, host, port, dbname)
	fmt.Println(scheme_url)
	db, err := gorm.Open("mysql", scheme_url)
	if err != nil {
		panic("failed to connect database mysql, error: " + err.Error())
	}
	// Migrate the schema
	return &DatabaseConnection{db}
}

func (db *DatabaseConnection) GetUrls() (bots []Bot) {
	db.Find(&bots)
	return bots;
}

func (db *DatabaseConnection) AddUrl(bot *Bot) {
	db.Create(bot)
}

func (db *DatabaseConnection) AddText(text *PasteText) {
	db.Create(text)
}

func (db *DatabaseConnection) GetText(code string) (text PasteText) {
	db.Where("code = ?", code).First(&text)
	return text;
}