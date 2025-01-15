package storage

import (
	"errors"
	"time"
)

var (
	ErrURLNotFound = errors.New("url not found")
	ErrURLExists   = errors.New("url already exists")
	ErrNoRows      = errors.New("no rows affected")
)

type Url struct {
	ID    int64     `gorm:"primaryKey"`
	Alias string    `gorm:"not null;unique;"`
	Url   string    `gorm:"not null"`
	Exp   time.Time `gorm:"not null"`
}
