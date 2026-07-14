package postgres

import (
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type chargePointStatus int

const (
	statusOffline chargePointStatus = iota - 1
	statusIdle
	statusOccupied
)

var StatusNames = map[chargePointStatus]string{
	statusOffline:  "offline",
	statusIdle:     "available",
	statusOccupied: "occupied",
}

func (s chargePointStatus) String() string {
	return StatusNames[s]
}

type ChargePoint struct {
	gorm.Model
	VendorId  string
	Name      string
	Latitude  float64
	Longitude float64
	Status    chargePointStatus
}

func (c *ChargePoint) SetStatus(statusString string) {
	for k, v := range StatusNames {
		if v == statusString {
			c.Status = k
			return
		}
	}
}
