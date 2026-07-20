package postgres

import (
	"testing"

	"github.com/restayway/gogis"
	"github.com/stretchr/testify/assert"
	"github.com/vpmv/chargepoint-api/internal/dto"
)

func TestMapper(t *testing.T) {

	dtoType := []*dto.ChargePoint{
		{
			VendorId: "Test1",
			Location: dto.Location{
				Latitude:  11,
				Longitude: 22,
			},
			Status: "available",
		},
	}
	postgresType := []*ChargePoint{
		{
			VendorId: "Test1",
			Point:    gogis.Point{Lat: 11, Lng: 22},
			Status:   statusIdle,
		},
	}

	assert.Equal(t, dtoType, chargePointsToDTOSlice(postgresType))
	assert.Equal(t, postgresType, dtoSliceToModel(dtoType))
}
