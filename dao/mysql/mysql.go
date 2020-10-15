package mysql

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/spf13/viper"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func Init() (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True",
		viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetString("mysql.port"),
		viper.GetString("mysql.db_name"),
	)

	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		zap.L().Error("connect to DB failed", zap.Error(err))
		return
	}

	err = db.Ping()
	if err != nil {
		fmt.Printf("connect to mysql failed, err: %v\n", err)
		return
	}

	db.SetMaxOpenConns(viper.GetInt("mysql.max_open_conns"))
	db.SetMaxIdleConns(viper.GetInt("mysql.max_idle_conns"))

	return nil
}

func Close() {
	_ = db.Close()
}
