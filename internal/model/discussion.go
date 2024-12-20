package model

import (
    "errors"
    "gorm.io/gorm"
)

type Discussion struct {
    ID      uint64      `gorm:"primaryKey;autoIncrement;index"`
    UserId  uint64      `gorm:"not null"`
    Title   string      `gorm:"not null"`
    Content string      `gorm:"not null"`
    Time    int64
}

var ErrDiscussionNotFound = errors.New("discussion not found")

func CreateDiscussion(discussion *Discussion) error {
    return db.Model(&Discussion{}).Create(discussion).Error
}

func SaveDiscussion(discussion *Discussion) error {
    return db.Save(discussion).Error
}

func FindDiscussionById(id uint64) (*Discussion, error) {
    discussion := &Discussion{}
    result := db.Model(&Discussion{}).Where("id = ?", id).First(discussion)
    if errors.Is(result.Error, gorm.ErrRecordNotFound) {
        return nil, ErrDiscussionNotFound
    }
    return discussion, result.Error
}