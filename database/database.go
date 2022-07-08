package database

import (
	"fmt"
	"sync"

	"github.com/alterra-dataon-kelompok-1/order-service/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	dbConn *gorm.DB
	once   sync.Once
)

func ConnectSQLDB(cfg config.Config) (*gorm.DB, error) {
	dbUser := cfg.Get("DB_USER")
	dbPassword := cfg.Get("DB_PASSWORD")
	dbUrl := cfg.Get("DB_URL")
	dbPort := cfg.Get("DB_PORT")
	dbName := cfg.Get("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPassword, dbUrl, dbPort, dbName)
	fmt.Println(dsn)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
}

// func CreateConnection(cfg config.Config) {
// 	dbUser := cfg.Get("DB_USER")
// 	dbPassword := cfg.Get("DB_PASSWORD")
// 	dbUrl := cfg.Get("DB_URL")
// 	dbPort := cfg.Get("DB_PORT")
// 	dbName := cfg.Get("DB_NAME")
//
// 	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPassword, dbUrl, dbPort, dbName)
// 	fmt.Println(dsn)
// 	once.Do(func() {
// 		dbConn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
// 		if err != nil {
// 			panic(err)
// 		}
// 	})
// }
//
// func GetConnection() *gorm.DB {
// 	if dbConn == nil {
// 		CreateConnection()
// 	}
// 	return dbConn
// }

// func ConnectDevDB() (*gorm.DB, error) {
// 	return gorm.Open(sqlite.Open("dev.db"), &gorm.Config{})
// }
