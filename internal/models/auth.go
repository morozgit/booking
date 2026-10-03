package models

type UserModel struct {
	ID             uint   `gorm:"primaryKey"`
	Email          string `gorm:"size:200;uniqueIndex"`
	HashedPassword string `gorm:"size:200"`
}

func (UserModel) TableName() string {
	return "users"
}
