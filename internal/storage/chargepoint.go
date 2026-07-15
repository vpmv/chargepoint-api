package storage

import "github.com/vpmv/chargepoint-api/internal/dto"

type ChargePointClient interface {
	Migrate() error
	GetChargePoints(page, pageSize int) (dto.ChargePoints, error)
	GetChargePoint(vendorID string) (*dto.ChargePoint, error)
	CreateChargePoints(chargePoint dto.ChargePoints) (dto.ChargePoints, error)
	DeleteChargePoint(vendorID string) error
	GetChargePointByLocation(latitude, longitude float64, radiusKm int) (dto.ChargePoints, error)
	MustSeed()
}
