package storage

import "github.com/vpmv/chargepoint-api/internal/dto"

type ChargePointClient interface {
	Migrate() error
	GetChargePoints(page, pageSize int) (*dto.ChargePoints, error)
	GetChargePoint(id string) (*dto.ChargePoint, error)
	CreateChargePoints(chargePoint *dto.ChargePoints) (*dto.ChargePoints, error)
	DeleteChargePoint(id string) error
	GetChargePointByLocation(latitude, longitude float64, radiusKm int) *dto.ChargePoints
}
