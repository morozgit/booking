package models

type HotelsModel struct {
	ID       uint   `gorm:"primaryKey"`
	Title    string `gorm:"size:100"`
	Location string `gorm:"size:400"`
}

func (HotelsModel) TableName() string {
	return "hotels"
}
