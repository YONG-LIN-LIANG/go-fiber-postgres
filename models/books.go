package models

import "gorm.io/gorm"

type Books struct{
	ID					uint					`gorm:"primary key" json:"id"`
	Author			*string				`json:"author"`
	Title				*string				`json:"title"`
	Publisher   *string				`json:"publisher"`
}

func MigrateBooks(db *gorm.DB) error{
	// AutoMigrate() 方法在 db 物件上被調用，傳遞 &Books{} 結構的指標作為參數。這個方法會自動在資料庫中建立 books 表格
	err := db.AutoMigrate(&Books{})
	return err
}