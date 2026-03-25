package ds

import (
	"time"
)

type Nuclear_calculation struct {
	CalcID      uint       `gorm:"primaryKey;column:calc_id"`
	Status      string     `gorm:"type:varchar(20);not null"`
	CreatedAt   time.Time  `gorm:"not null"`
	CreatorID   uint       `gorm:"not null"`
	FormingDate *time.Time `gorm:"column:forming_date"`
	FinishDate  *time.Time `gorm:"column:finish_date"`
	ModeratorID *uint      `gorm:"column:moderator_id"`
	Description *string    `gorm:"type:varchar(2000)"`

	Creator   Engineers  `gorm:"foreignKey:CreatorID"`
	Moderator *Engineers `gorm:"foreignKey:ModeratorID"`
}

func (Nuclear_calculation) TableName() string {
	return "nuclear_calculations"
}
