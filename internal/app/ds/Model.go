package ds

type Model struct {
	ModelID     uint    `gorm:"primaryKey;column:model_id"`
	Title       string  `gorm:"type:varchar(255);not null"`
	Description string  `gorm:"type:varchar(1000);"`
	ShortDesc   string  `gorm:"type:varchar(100);`
	IsDeleted   bool    `gorm:"type:boolean;not null;default:false"`
	PhotoURL    string  `gorm:"column:photo_url;type:varchar(255)"`
	Video       string  `gorm:"type:varchar(255)"`
	Power       float64 `gorm:"type:numeric(6,2);not null"`
	FuelUsage   float64 `gorm:"type:numeric(6,2);not null"`
}

func (Model) TableName() string {
	return "models"
}
