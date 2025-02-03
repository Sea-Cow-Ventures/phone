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
	Cnf    config.AppConfig
	db     *sqlx.DB
	Logger *zap.SugaredLogger
)

func init() {
	Cnf = config.GetConfig()
	Logger = log.GetLogger()
}

func InitMysql() {
	mCnf := mysql.Config{
		User:                 Cnf.DBUser,
		Passwd:               Cnf.DBPass,
		Net:                  "tcp",
		Addr:                 Cnf.DBServer + ":3306",
		DBName:               Cnf.DBSchema,
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	mysql, err := sqlx.Open("mysql", mCnf.FormatDSN())

	mysql.SetConnMaxLifetime(time.Minute * 3)
	mysql.SetMaxOpenConns(10)
	mysql.SetMaxIdleConns(10)

	if err != nil {
		Logger.Fatal("Failed to connect to mysql db", zap.Error(err))
	}

	err = mysql.Ping()
	if err != nil {
		Logger.Fatal("Failed to connect to mysql db", zap.Error(err))
	}

	db = mysql
	Logger.Info("Connected to db")
}

func GetDb() *sqlx.DB {
	return db
}
