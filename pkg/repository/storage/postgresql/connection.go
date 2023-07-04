package postgresql

import (
	"fmt"
	"net/url"
	"os"
	"strconv"

	"log"

	"github.com/vesicash/mor-api/internal/config"
	"github.com/vesicash/mor-api/utility"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	lg "gorm.io/gorm/logger"
)

type Databases struct {
	Admin         *gorm.DB
	Auth          *gorm.DB
	Notifications *gorm.DB
	Payment       *gorm.DB
	Reminder      *gorm.DB
	Subscription  *gorm.DB
	Transaction   *gorm.DB
	Verification  *gorm.DB
	Cron          *gorm.DB
	MOR           *gorm.DB
}

var DB Databases

// Connection gets connection of mysqlDB database
func Connection() Databases {
	return DB
}

func ConnectToDatabases(logger *utility.Logger, configDatabases config.Databases) Databases {
	dbsCV := configDatabases
	databases := Databases{}
	utility.LogAndPrint(logger, "connecting to databases")
	databases.MOR = connectToDb(dbsCV.DB_HOST, dbsCV.USERNAME, dbsCV.PASSWORD, dbsCV.MOR_DB, dbsCV.DB_PORT, dbsCV.SSLMODE, dbsCV.TIMEZONE, logger)

	utility.LogAndPrint(logger, "connected to databases")

	utility.LogAndPrint(logger, "connected to db")
	// migrations

	DB = databases
	return DB
}

func connectToDb(host, user, password, dbname, port, sslmode, timezone string, logger *utility.Logger) *gorm.DB {
	if _, err := strconv.Atoi(port); err != nil {
		u, err := url.Parse(port)
		if err != nil {
			utility.LogAndPrint(logger, fmt.Sprintf("parsing url %v to get port failed with: %v", port, err))
			panic(err)
		}

		detectedPort := u.Port()
		if detectedPort == "" {
			utility.LogAndPrint(logger, fmt.Sprintf("detecting port from url %v failed with: %v", port, err))
			panic(err)
		}
		port = detectedPort
	}
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=%v TimeZone=%v", host, user, password, dbname, port, sslmode, timezone)

	newLogger := lg.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		lg.Config{
			LogLevel:                  lg.Error, // Log level
			IgnoreRecordNotFoundError: true,     // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,
		},
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		utility.LogAndPrint(logger, fmt.Sprintf("connection to %v db failed with: %v", dbname, err))
		panic(err)

	}

	utility.LogAndPrint(logger, fmt.Sprintf("connected to %v db", dbname))
	return db
}

func ReturnDatabase(name string) *gorm.DB {
	databases := DB
	switch name {
	case "admin_service":
		return DB.Admin
	case "auth_service":
		return DB.Auth
	case "notification_service":
		return DB.Notifications
	case "payment_service":
		return DB.Payment
	case "reminders_service":
		return DB.Reminder
	case "subscription_service":
		return DB.Subscription
	case "transaction_service":
		return DB.Transaction
	case "verification_service":
		return DB.Verification
	case "cron_service":
		return DB.Cron
	case "mor":
		return DB.MOR
	default:
		return databases.Auth
	}
}
