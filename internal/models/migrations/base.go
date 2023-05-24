package migrations

import (
	"github.com/vesicash/mor-api/pkg/repository/storage/postgresql"
	"gorm.io/gorm"
)

func RunAllMigrations(db postgresql.Databases) {

	// payment migration
	MigrateModels(db.MOR, AuthMigrationModels())

}

func MigrateModels(db *gorm.DB, models []interface{}) {
	_ = db.AutoMigrate(models...)
}
