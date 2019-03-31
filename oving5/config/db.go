package config

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql" // MySql driver used by gorm
	"github.com/jinzhu/gorm"
)

const depth = 1

// DB is the database instance
var DB *gorm.DB

func init() {
	fmt.Println("Connecting to database")
	db, err := gorm.Open("mysql", "sveinuov:4MHEFBGw@tcp(mysql.stud.iie.ntnu.no)/sveinuov?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	db.DB().SetMaxOpenConns(depth)
	// db.LogMode(true)

	DB = db
}
