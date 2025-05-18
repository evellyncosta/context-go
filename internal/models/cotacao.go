package models

import "gorm.io/gorm"

type Cotacao struct {
	gorm.Model
	Valor    float64 `gorm:"not null"`
	Currency string  `gorm:"not null"`
}