package model

import (
    "fmt"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"

    "user-management-system/internal/config"
)

var db *gorm.DB

func InitDB() {
    dsn := config.Config.Server.Dsn
    DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        panic(err)
    }
    db = DB
    if err := db.AutoMigrate(&User{}); err != nil {
        panic(err)
    }
    if err := db.AutoMigrate(&Discussion{}); err != nil {
        panic(err)
    }
    if err := db.AutoMigrate(&Comment{}); err != nil {
        panic(err)
    }
}