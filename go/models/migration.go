package models

import "gorm.io/gorm"

func DatabaseMigration(db *gorm.DB){
	db.AutoMigrate(User{})
	db.AutoMigrate(Token{})
}