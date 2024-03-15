package database

import (
	"fmt"

	"github.com/maxuanquang/idm/internal/configs"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitializeDB(dbConfig configs.Database) (*gorm.DB, func(), error) {
	// Create data source name (DSN) string
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", dbConfig.Username, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Database)

	// Open GORM database connection
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}

	return db, cleanup, nil
}
