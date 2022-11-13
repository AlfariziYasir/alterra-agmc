package model

import (
	"time"

	"gorm.io/gorm"
)

type Base struct {
	CreatedAt time.Time       `gorm:"column:created_at"`
	CreatedBy string          `json:"created_by" gorm:"type:varchar(25);default:Admin"`
	UpdatedAt time.Time       `gorm:"column:updated_at"`
	UpdatedBy string          `json:"updated_by" gorm:"type:varchar(25);default:Admin"`
	DeletedAt *gorm.DeletedAt `gorm:"index"`
	DeletedBy string          `json:"deleted_by" gorm:"type:varchar(25);"`
	ID        uint            `gorm:"primaryKey;index;NOT NULL;column:id;autoIncrement"`
}

func (b *Base) BeforeCreate(tx *gorm.DB) (err error) {
	now := time.Now()
	b.CreatedAt = now
	b.UpdatedAt = now
	return
}

func (b *Base) BeforeUpdate(tx *gorm.DB) (err error) {
	b.UpdatedAt = time.Now()
	return
}

func (b *Base) BeforeDelete(tx *gorm.DB) (err error) {
	b.DeletedAt = &gorm.DeletedAt{Time: time.Now(), Valid: true}
	return
}
