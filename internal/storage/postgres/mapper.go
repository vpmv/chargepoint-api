package postgres

import (
	"github.com/restayway/gogis"
	"github.com/vpmv/chargepoint-api/internal/dto"
	"github.com/vpmv/chargepoint-api/internal/storage"
)

func dtoSliceToModel(dto []*dto.ChargePoint) []*ChargePoint {
	return storage.MapSlice(dto, dtoToChargePoint)
}

func chargePointsToDTOSlice(model []*ChargePoint) []*dto.ChargePoint {
	return storage.MapSlice(model, chargePointToDTO)
}

func dtoToChargePoint(dto *dto.ChargePoint) *ChargePoint {
	c := &ChargePoint{
		VendorId: dto.VendorId,
		Name:     dto.Name,
		Point:    gogis.Point{Lat: dto.Location.Latitude, Lng: dto.Location.Longitude},
	}
	c.SetStatus(dto.Status)

	return c
}

func chargePointToDTO(chargePoint *ChargePoint) *dto.ChargePoint {
	return &dto.ChargePoint{
		VendorId: chargePoint.VendorId,
		Name:     chargePoint.Name,
		Location: dto.Location{
			Latitude:  chargePoint.Point.Lat,
			Longitude: chargePoint.Point.Lng,
		},
		Status:   chargePoint.Status.String(),
		Distance: chargePoint.Distance / 1000,
	}
}
