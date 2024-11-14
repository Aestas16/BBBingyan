package model

type Comment struct {
    ID      uint64      `gorm:"primaryKey;autoIncrement;index"`
    UserId  uint64      `gorm:"not null"`
    DiscId  uint64      `gorm:"not null;index"`
    Content string      `gorm:"not null"`
    Time    int64
}

func CreateComment(comment *Comment) error {
    return db.Model(&Comment{}).Create(comment).Error
}

func FindCommentsByDiscId(id uint64) ([]Comment, error) {
    var comments []Comment
    result := db.Model(&Comment{}).Where("discid = ?", id).Find(&comments)
    return comments, result.Error
}