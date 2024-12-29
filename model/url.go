package model

type Url struct {
	ID    int64  `gorm:"primaryKey"`
	Alias string `gorm:"not null;unique;"`
	Url   string `gorm:"not null"`
}
