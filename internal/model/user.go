package model

import (
    "errors"
    "gorm.io/gorm"
)

type User struct {
    ID          uint64  `gorm:"primaryKey;autoIncrement;index"`
    Username    string  `gorm:"type:varchar(80);unique;not null;" json:"username"`
    Password    string  `json:"password"`
    Email       string  `json:"email"`
}

var ErrUserNotFound = errors.New("user not found")
var ErrUserAlreadyExist = errors.New("user already exist")

func CreateUser(user *User) error {
    _, err := FindUserByName(user.Username)
    if err == nil {
        return ErrUserAlreadyExist
    }
    return db.Model(&User{}).Create(user).Error
}

func SaveUser(user *User) error {
    return db.Save(user).Error
}

func DeleteUserById(id uint64) error {
    _, err := FindUserById(id)
    if err != nil {
        return err
    }
    return db.Delete(&User{}, id).Error
}

func FindUserByName(username string) (*User, error) {
    user := &User{}
    result := db.Where("username = ?", username).First(user)
    if errors.Is(result.Error, gorm.ErrRecordNotFound) {
        return nil, ErrUserNotFound
    }
    return user, result.Error
}

func FindUserById(id uint64) (*User, error) {
    user := &User{}
    result := db.Model(&User{}).Where("id = ?", id).First(user)
    if errors.Is(result.Error, gorm.ErrRecordNotFound) {
        return nil, ErrUserNotFound
    }
    return user, result.Error
}