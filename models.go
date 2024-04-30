package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"gorm.io/gorm"
)

var validate = validator.New()
var bundle *i18n.Bundle
var db *gorm.DB

// GORM model for FAQ data
type FAQ struct {
	ID         uuid.UUID `gorm:"primaryKey"`
	Question   string
	TrQuestion string
	Answers    []Answer `gorm:"foreignKey:FAQID"`
}

type Answer struct {
	ID       *uint     `gorm:"primaryKey"`
	FAQID    uuid.UUID `json:"-"`
	Title    string
	TrTitle  string
	Answer   string
	TrAnswer string
}
