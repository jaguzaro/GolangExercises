package models

type Users struct {
	ID       int    `gorm:"unique;primaryKey;autoIncrement;column:id;"`
	Name     string `gorm:"column:name"`
	Gender   string `gorm:"column:gender"`
	Location string `gorm:"column:location"`
	Email    string `gorm:"column:email"`
	Phone    string `gorm:"column:phone"`
}
