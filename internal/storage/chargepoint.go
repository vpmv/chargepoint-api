package storage

import "github.com/vpmv/chargepoint-api/internal/dto"

type ChargePointClient interface {
	Migrate() error
	GetChargePoints() ([]*dto.ChargePoint, error)
	GetChargePoint(id string) (*dto.ChargePoint, error)
	CreateChargePoint(chargePoint *dto.ChargePoint) (*dto.ChargePoint, error)
	DeleteChargePoint(id string) error
	GetChargePointByLocation(latitude, longitude float64, radiusKm int) []*dto.ChargePoint
}
