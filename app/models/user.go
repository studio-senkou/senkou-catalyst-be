package models

import (
	"time"

	"gorm.io/gorm"

	bcrypt "golang.org/x/crypto/bcrypt"
)

type User struct {
	ID              uint32         `json:"id"           gorm:"type:int;primaryKey"`
	Merchants       []*Merchant    `json:"merchants"    gorm:"foreignKey:OwnerID;references:ID"`
	Name            string         `json:"name"         gorm:"type:varchar(100);not null"`
	Email           string         `json:"email"        gorm:"type:varchar(100);unique;not null"`
	Phone           string         `json:"phone"        gorm:"type:varchar(20);unique;not null"`
	Password        []byte         `json:"-"            gorm:"type:varchar(255);not null"`
	Role            string         `json:"role"         gorm:"type:varchar(20);not null;default:user"`
	EmailVerifiedAt *time.Time     `json:"email_verified_at" gorm:"type:timestamp;default:null"`
	CreatedAt       time.Time      `json:"created_at"   gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt       time.Time      `json:"updated_at"   gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	DeletedAt       gorm.DeletedAt `json:"-"            gorm:"type:timestamp;index"`
}

func (u *User) HashPassword() ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
}

func (u *User) CheckPassword(password string) bool {
	return bcrypt.CompareHashAndPassword(u.Password, []byte(password)) == nil
}
