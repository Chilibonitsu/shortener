package model

import "time"

type Url struct {
	ID    int64     `gorm:"primaryKey"`
	Alias string    `gorm:"not null;unique;"`
	Url   string    `gorm:"not null"`
	Exp   time.Time `gorm:"not null"`
	//Если вдруг нужно добавить поле, то что обычно делают?
}
