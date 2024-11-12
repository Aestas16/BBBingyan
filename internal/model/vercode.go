package model

import (
    "errors"
    "gorm.io/gorm"
	"time"
)

type VerCode struct {
	UserId	uint64  `gorm:"primaryKey;index"`
	Code	string	`gorm:"type:varchar(6);not null;"`
	Time	time.Time
}

var ErrVerCodeNotFound = errors.New("verification code not found")
var ErrVerCodeAlreadyExist = errors.New("verification code already exist")

func CreateVerCode(vercode *VerCode) error {
    _, err := FindVerCodeByUserId(vercode.UserId)
    if err == nil {
        return ErrVerCodeAlreadyExist
    }
    return db.Model(&VerCode{}).Create(vercode).Error
}

func SaveVerCode(vercode *VerCode) error {
	return db.save(vercode).Error
}

func FindVerCodeByUserId(uid uint64) (*VerCode, error) {
    vercode := &VerCode{}
    result := db.Model(&VerCode{}).Where("userid = ?", uid).First(vercode)
    if errors.Is(result.Error, gorm.ErrRecordNotFound) {
        return nil, ErrVerCodeNotFound
    }
    return vercode, result.Error
}