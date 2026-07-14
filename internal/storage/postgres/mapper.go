package postgres

import (
	"strconv"

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
		VendorId:  dto.VendorId,
		Name:      dto.Name,
		Latitude:  dto.Location.Latitude,
		Longitude: dto.Location.Longitude,
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
			Latitude:  chargePoint.Latitude,
			Longitude: chargePoint.Longitude,
		},
		Status: chargePoint.Status.String(),
	}
}
