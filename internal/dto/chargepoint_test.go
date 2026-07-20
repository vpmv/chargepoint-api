package dto

import (
	"context"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

var _testChargePoints = []*ChargePoint{
	{
		VendorId: "SH001-001",
		Name:     "Shell A2 - Mock 001",
		Location: Location{
			Latitude:  3.14,
			Longitude: 31.1,
		},
		Status: "AVAILABLE",
	},
	{
		VendorId: "SH001-002",
		Name:     "Shell A2 - Mock 002",
		Location: Location{
			Latitude:  3.141,
			Longitude: 31.1,
		},
		Status: "OCCUPIED",
	},
	{
		VendorId: "SH001-003",
		Name:     "Shell A2 - Mock 003",
		Location: Location{
			Latitude:  3.142,
			Longitude: 31.1,
		},
	},
	{
		VendorId: "SH001-004",
		Name:     "Shell A2 - Mock 004",
		Location: Location{
			Latitude:  3.143,
			Longitude: 31.1,
		},
		Status: "CORRUPT",
	},
	{
		VendorId: "SH001-005",
		Name:     "Shell A2 - Mock 005",
		Location: Location{
			Latitude:  0,
			Longitude: 0,
		},
		Status: "OFFLINE",
	},
}

func _copyChargePoints() ChargePoints {
	out := make([]*ChargePoint, len(_testChargePoints))
	for i := range _testChargePoints {
		cp := *_testChargePoints[i]
		out[i] = new(cp)
	}
	return out
}

func TestChargePoints_InTransform(t *testing.T) {
	var chargePoints = _copyChargePoints()

	if err := chargePoints.InTransform(context.Background()); err != nil {
		t.Error(`error performing ChargePoints.InTransform()`)
	}

	assert.Equal(t, "available", chargePoints[0].Status, `status should be lower-cased`)
	assert.Equal(t, "available", chargePoints[2].Status, `status should be set to default`)
	assert.Equal(t, "occupied", chargePoints[1].Status, `status should be lower-cased`)
	assert.Equal(t, "offline", chargePoints[4].Status, `status should be lower-cased`)
}

func TestChargePoints_Validate(t *testing.T) {
	var (
		chargePoints = _copyChargePoints()
		errs         map[string]validator.ValidationErrors
	)

	chargePoints, errs = chargePoints.Validate(true)
	assert.Equal(t, 2, len(chargePoints), `Three ChargePoints should have been dismissed`)
	assert.Equal(t, 3, len(errs), `Three ChargePoints should be faulty`)

	assert.Equal(t, `Status`, errs[`SH001-003`][0].Field())
	assert.Equal(t, `Status`, errs[`SH001-004`][0].Field())
	assert.Equal(t, `Latitude`, errs[`SH001-005`][0].Field())

	chargePoints = _copyChargePoints()
	_ = chargePoints.InTransform(context.Background())
	chargePoints, errs = chargePoints.Validate(false)
	assert.Zero(t, len(chargePoints), `ChargePoints should have been reset`)
	assert.Equal(t, 1, len(errs), `One ValidationError should have been created`)

}
