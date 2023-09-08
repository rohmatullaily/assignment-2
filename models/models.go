package models

type User struct {
	ID        uint   `gorm:"primaryKey;type:int;autoIncrement"`
	Username  string `gorm:"column:username;type:varchar(50);not null"`
	FirstName string `gorm:"column:first_name;type:varchar(200);not null"`
	LastName  string `gorm:"column:last_name;type:varchar(200);not null"`
	Password  string `gorm:"column:password;type:varchar(120);not null"`
}

type Layout struct{
	Message  string
}