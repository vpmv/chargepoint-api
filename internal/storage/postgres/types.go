package postgres

import (
	"github.com/restayway/gogis"
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
	Name     string
	Status   chargePointStatus
	VendorId string      `gorm:"uniqueIndex"`
	Point    gogis.Point `gorm:"type:geography(Point, 4326)"`
	Distance float64     `gorm:"column:distance;->"` // virtual field, only available for ByLocation query
}

func (c *ChargePoint) SetStatus(statusString string) {
	for k, v := range StatusNames {
		if v == statusString {
			c.Status = k
			return
		}
	}
}
