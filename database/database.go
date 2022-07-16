package database

import (
	"fmt"
	"log"
  "pro1/entity"
	"github.com/jinzhu/gorm"
	_"github.com/jinzhu/gorm/dialects/postgres"
)


func GetDataBase() *gorm.DB {
	dbname := "covid"
	db := "postgres"
	db_password := "123456"

	db_url := "postgres://postgres:" + db_password + "@localhost/" + dbname + "?sslmode=disable"

	conn, err := gorm.Open(db, db_url)
	if err != nil {
		log.Fatalln("Invalid database url")
	}

	sqldb := conn.DB()
	err = sqldb.Ping()
	if err != nil {
		log.Fatalln("Database connected!")
	}
	fmt.Println("Database connection successful.")

	conn.AutoMigrate(&entity.Data{}, &entity.User{})

	return conn

}

func CloseDataBase(conn *gorm.DB) {
	sql := conn.DB()
	sql.Close()
}
