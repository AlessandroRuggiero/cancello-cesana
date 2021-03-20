package database

import (
	"fmt"
	"time"
)

func (d *CesanaDb) AddAperturaCancello(id string) string {
	var apertura Apertura
	d.Db.First(&apertura, id)
	if apertura.ID == 0 {
		fmt.Println("No apertura in db")
		return ""
	}
	apertura.Done = true
	d.Db.Save(&apertura)
	return apertura.AuthId
}

func (d *CesanaDb) GetApertureCancello(id string) int {
	var aperture []Apertura
	ieri := time.Now().Add(-24 * time.Hour)
	d.Db.Where("type = ? AND created_at > ? AND auth_id = ? AND done = true", "cancello", ieri, id).Find(&aperture)
	return len(aperture)
}

func (d *CesanaDb) AddAperturaCancelloAttempt(id string) uint {
	apertura := &Apertura{
		AuthId: id,
		Type:   "cancello",
	}
	d.Db.Create(apertura)
	fmt.Println(apertura.ID)
	return apertura.ID
}

func (d *CesanaDb) GetApertureCancelletto(id string) int {
	var aperture []Apertura
	ieri := time.Now().Add(-24 * time.Hour)
	d.Db.Where("type = ? AND created_at > ? AND auth_id = ? AND done = true", "cancelletto", ieri, id).Find(&aperture)
	return len(aperture)
}

func (d *CesanaDb) AddAperturaCancellettoAttempt(id string) uint {
	apertura := &Apertura{
		AuthId: id,
		Type:   "cancelletto",
	}
	d.Db.Create(apertura)
	fmt.Println(apertura.ID)
	return apertura.ID
}
func (d *CesanaDb) AddAperturaCancelletto(id string) string {
	fmt.Println("Id nell apertura: ", id)
	var apertura Apertura
	d.Db.First(&apertura, id)
	if apertura.ID == 0 {
		fmt.Println("impossibile trovare log")
		return ""
	}
	apertura.Done = true
	d.Db.Save(&apertura)
	return apertura.AuthId
}

func (d *CesanaDb) GetAperture(id string) *Aperture {
	return &Aperture{
		Cancello:    d.GetApertureCancello(id),
		Cancelletto: d.GetApertureCancelletto(id),
	}
}
