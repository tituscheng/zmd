package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type SQLite struct {
	Path string
	db   *gorm.DB
}

func (sql SQLite) Connect() *gorm.DB {
	var err error
	if sql.db == nil {
		sql.db, err = gorm.Open(sqlite.Open(sql.Path), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err != nil {
			panic("failed to connect database")
		}
	}
	return sql.db
}
