package model

import (
    "errors"
    "gorm.io/gorm"
	"time"
)

type Comment struct {
    ID		uint64  	`gorm:"primaryKey;autoIncrement;index"`
    UserId	uint64		`gorm:"not null"`
	DiscId	uint64		`gorm:"not null;index"`
	Content	string  	`gorm:"not null"`
    Time	time.Time
}

func CreateComment(discussion *Discussion) error {
    return db.Model(&Discussion{}).Create(discussion).Error
}

func FindCommentsByDiscId(id uint64) ([]Comment, error) {
    var comments []Comment
    result := db.Model(&Comment{}).Where("discid = ?", id).Find(&comments)
    return comments, result.Error
}