package database

import (
	"aidan/phone/internal/config"
	"aidan/phone/internal/log"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var (
	db *sqlx.DB
)

func init() {
	cnf := config.GetConfig()
	logger := log.GetLogger()

	mCnf := mysql.Config{
		User:                 cnf.DBUser,
		Passwd:               cnf.DBPass,
		Net:                  "tcp",
		Addr:                 cnf.DBServer + ":3306",
		DBName:               cnf.DBSchema,
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	mysql, err := sqlx.Open("mysql", mCnf.FormatDSN())

	mysql.SetConnMaxLifetime(time.Minute * 3)
	mysql.SetMaxOpenConns(10)
	mysql.SetMaxIdleConns(10)

	if err != nil {
		logger.Fatal("Failed to connect to mysql db", zap.Error(err))
	}

	err = mysql.Ping()
	if err != nil {
		logger.Fatal("Failed to connect to mysql db", zap.Error(err))
	}

	db = mysql
	logger.Info("Connected to db")
}

func GetDb() *sqlx.DB {
	return db
}
