package ds

type Engineers struct {
	EngineerID   uint   `gorm:"primaryKey;column:engineer_id"`
	Login        string `gorm:"type:varchar(50);unique;not null"`
	Password     string `gorm:"type:varchar(100);not null"`
	IsTechnician bool   `gorm:"type:boolean;default:false"`
}

func (Engineers) TableName() string {
	return "engineers"
}
