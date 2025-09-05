package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID           string          `json:"id" gorm:"type:uuid;primaryKey"`
	MerchantID   string          `json:"merchant_id" gorm:"type:char(16);not null"`
	Merchant     Merchant        `json:"-" gorm:"foreignKey:MerchantID;references:ID"`
	CategoryID   *uint32         `json:"category_id" gorm:"type:int;default:null"`
	Category     *Category       `json:"-" gorm:"foreignKey:CategoryID;references:ID"`
	Interactions []ProductMetric `json:"-" gorm:"foreignKey:ProductID;references:ID"`
	Title        string          `json:"title" gorm:"type:varchar(150);not null"`
	Price        string          `json:"price" gorm:"type:varchar(30);not null"`
	Description  string          `json:"description" gorm:"type:text"`
	AffiliateURL string          `json:"affiliate_url" gorm:"type:text;not null"`
	Photos       PhotoArray      `json:"photos" gorm:"type:json;default:'[]'"`
	CreatedAt    time.Time       `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time       `json:"updated_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	DeletedAt    gorm.DeletedAt  `json:"-" gorm:"type:timestamp;index"`
}

type PhotoArray []string

func (p PhotoArray) Value() (driver.Value, error) {
	if len(p) == 0 {
		return "[]", nil
	}
	return json.Marshal(p)
}

func (p *PhotoArray) Scan(value any) error {
	if value == nil {
		*p = PhotoArray{}
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New("cannot scan info photo array")
	}

	return json.Unmarshal(bytes, p)
}

func (p PhotoArray) MarshalJSON() ([]byte, error) {
	return json.Marshal([]string(p))
}

func (p *PhotoArray) UnmarshalJSON(data []byte) error {
	var strings []string
	if err := json.Unmarshal(data, &strings); err != nil {
		return err
	}

	*p = PhotoArray(strings)
	return nil
}

func (p *PhotoArray) AddPhoto(photo string) {
	*p = append(*p, photo)
}

func (p *PhotoArray) RemovePhoto(photo string) {
	for i, v := range *p {
		if v == photo {
			*p = append((*p)[:i], (*p)[i+1:]...)
			break
		}
	}
}
