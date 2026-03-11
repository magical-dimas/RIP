package ds

type ModelRequest struct {
	RequestID uint `gorm:"primaryKey;column:request_id"`
	ModelID   uint `gorm:"primaryKey;column:model_id"`

	Amount   int      `gorm:"not null;default:1"`
	ResPower *float64 `gorm:"type:numeric(8,2)"`
	ResFuel  *float64 `gorm:"type:numeric(8,2)"`

	Model   Model   `gorm:"foreignKey:ModelID"`
	Request Request `gorm:"foreignKey:RequestID"`
}

func (ModelRequest) TableName() string {
	return "modelrequests"
}
