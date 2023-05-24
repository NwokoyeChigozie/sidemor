package migrations

import "github.com/vesicash/mor-api/internal/models"

// _ = db.AutoMigrate(MigrationModels()...)
func AuthMigrationModels() []interface{} {
	return []interface{}{
		models.Customer{},
		models.PaymentHistory{},
		models.PaymentModule{},
		models.PaymentOrder{},
		models.Setting{},
	}
}
