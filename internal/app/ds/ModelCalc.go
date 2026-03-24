package ds

type ModelCalc struct {
	CalcID  uint `gorm:"primaryKey;column:calc_id"`
	ModelID uint `gorm:"primaryKey;column:model_id"`

	Amount   int      `gorm:"not null;default:1"`
	ResPower *float64 `gorm:"type:numeric(8,2)"`
	ResFuel  *float64 `gorm:"type:numeric(8,2)"`

	Model               Model               `gorm:"foreignKey:ModelID"`
	Nuclear_calculation Nuclear_calculation `gorm:"foreignKey:CalcID"`
}

func (ModelCalc) TableName() string {
	return "modelcalcs"
}
