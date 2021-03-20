package database

import "github.com/jinzhu/gorm"

type CesanaDb struct {
	Db    *gorm.DB
	Ready bool
}
