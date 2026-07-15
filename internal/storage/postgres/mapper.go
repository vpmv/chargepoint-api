package postgres

import (
	"strconv"

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

	// parse ID if set, to enable updates
	if dto.ID != "" {
		id, err := strconv.ParseUint(dto.ID, 10, 64)
		if err == nil {
			c.ID = uint(id)
		}
	}
	return c
}

func chargePointToDTO(chargePoint *ChargePoint) *dto.ChargePoint {
	return &dto.ChargePoint{
		ID:       strconv.FormatUint(uint64(chargePoint.ID), 10),
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
