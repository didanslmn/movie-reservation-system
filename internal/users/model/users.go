package model

import (
	"time"

	"gorm.io/gorm"
)

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

type User struct {
	gorm.Model
	Name      string `gorm:"size:255;not null"`
	Email     string `gorm:"size:255;unique;not null"`
	Password  string `gorm:"size:255;not null"`
	Role      Role   `gorm:"type:role;default:'user';index"`
	LastLogin time.Time
}

func (Role) GormDataType() string {
	return "role"
}
