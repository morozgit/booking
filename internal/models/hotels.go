package models

type HotelsModel struct {
	ID       uint   `gorm:"primaryKey"`
	title    string `gorm:"size:100"`
	location string `gorm:"size:400"`
}

func (HotelsModel) TableName() string {
	return "hotels"
}
