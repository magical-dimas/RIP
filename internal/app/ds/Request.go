package ds

import (
	"database/sql"
	"time"
)

type Request struct {
	RequestID   uint         `gorm:"primaryKey;column:request_id"`
	Status      string       `gorm:"type:varchar(20);not null"`
	CreatedAt   time.Time    `gorm:"not null"`
	CreatorID   uint         `gorm:"not null"`
	FormingDate *time.Time   `gorm:"column:forming_date"`
	FinishDate  sql.NullTime `gorm:"column:finish_date"`
	ModeratorID *uint        `gorm:"column:moderator_id"`
	Description *string      `gorm:"type:varchar(2000)"`

	Creator   Users  `gorm:"foreignKey:CreatorID"`
	Moderator *Users `gorm:"foreignKey:ModeratorID"`
}

func (Request) TableName() string {
	return "requests"
}
