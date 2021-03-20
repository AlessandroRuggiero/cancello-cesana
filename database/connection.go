package database

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

//Connect connects to db
func (d *CesanaDb) Connect(dbURL string) error {
	db, err := gorm.Open("postgres", dbURL)
	if err != nil {
		return err
	}
	d.Db = db
	d.Ready = true
	d.Migrate()
	return nil
}

//Migrate the db
func (d *CesanaDb) Migrate() {
	d.Db.AutoMigrate(&Apertura{})
}

//Disconnect terminates connection with the server
func (d *CesanaDb) Disconnect() {
	fmt.Println("Disconnetto db")
	d.Db.Close()
}
