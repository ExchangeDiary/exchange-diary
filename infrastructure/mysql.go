package infrastructure

import (
	"fmt"

	"github.com/ExchangeDiary/exchange-diary/infrastructure/configs"
	"github.com/ExchangeDiary/exchange-diary/infrastructure/persistence"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func devDsn(cfg *configs.DBConfig) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)
}

func sandboxDsn(cfg *configs.DBConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
}

// ConnectDatabase returns
// https://gorm.io/docs/connecting_to_the_database.html
func ConnectDatabase(phase string) *gorm.DB {
	var dsn string
	if phase == "sandbox" {
		dsn = sandboxDsn(configs.DatabaseConfig())
	} else {
		dsn = devDsn(configs.DatabaseConfig())
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		panic(err.Error)
	}

	sqlDB, _ := db.DB()
	err = sqlDB.Ping()
	if err != nil {
		panic(err)
	}

	db.Set("gorm:table_options", "ENGINE=InnoDB")
	return db
}

// Migrate do db migrations
func Migrate(db *gorm.DB) {
	db.AutoMigrate(&persistence.AccountGorm{})
	db.AutoMigrate(&persistence.RoomGorm{})
	db.AutoMigrate(&persistence.RoomMemberGorm{})
	// db.AutoMigrate()
}
