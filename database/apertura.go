package database

import "github.com/jinzhu/gorm"

type Apertura struct {
	gorm.Model
	AuthId string `gorm:"index"`
	Type   string
	Done   bool
}

type Aperture struct {
	Cancello    int `json:"cancello"`
	Cancelletto int `josn:"cancelletto"`
}
