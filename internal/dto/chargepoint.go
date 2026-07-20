package dto

import (
	"context"
	"strings"

	"github.com/go-fuego/fuego"
	"github.com/go-playground/validator/v10"
	"github.com/vpmv/chargepoint-api/internal/helpers"
)

type Location struct {
	Latitude  float64 `json:"latitude" validate:"required,min=-90,max=90"`
	Longitude float64 `json:"longitude" validate:"required,min=-180,max=180"`
}

type ChargePoint struct {
	VendorId string   `json:"id" validate:"required"`
	Name     string   `json:"name" validate:"required"`
	Location Location `json:"location" validate:"required"`
	Status   string   `json:"status" validate:"oneof=OFFLINE AVAILABLE OCCUPIED"`
	Distance float64  `json:"distance,omitempty"` // output only
}

// InTransform - called by Fuego during ingress
func (c *ChargePoint) InTransform(ctx context.Context) error {
	c.Status = strings.ToLower(c.Status)
	if c.Status == `` {
		c.Status = `available`
	}
	return nil
}

// OutTransform - called by Fuego during egress
func (c *ChargePoint) OutTransform(ctx context.Context) error {
	c.Status = strings.ToUpper(c.Status)
	return nil
}

type ChargePoints []*ChargePoint

// InTransform - called by Fuego during ingress
func (s ChargePoints) InTransform(ctx context.Context) error {
	for i := range s {
		_ = s[i].InTransform(ctx)
	}
	return nil
}

// OutTransform - called by Fuego during egress
func (s ChargePoints) OutTransform(ctx context.Context) error {
	for i := range s {
		_ = s[i].OutTransform(ctx)
	}
	return nil
}

// Validate ingress chargepoints.
//
// Returns valid chargepoints and validation errors
func (s ChargePoints) Validate(ignoreErrors bool) (chargePoints ChargePoints, errs map[string]validator.ValidationErrors) {
	chargePoints = make([]*ChargePoint, 0)
	errs = make(map[string]validator.ValidationErrors)

	for _, cp := range s {
		validationErrors := helpers.ValidateEntity(cp)
		if validationErrors != nil {
			errs[cp.VendorId] = validationErrors
			if !ignoreErrors {
				return nil, errs
			}
		} else {
			chargePoints = append(chargePoints, cp)
		}
	}

	return chargePoints, errs
}

// type safety checks
var _ fuego.InTransformer = (*ChargePoint)(nil)
var _ fuego.OutTransformer = (*ChargePoint)(nil)

var _ fuego.InTransformer = (*ChargePoints)(nil)
var _ fuego.OutTransformer = (*ChargePoints)(nil)
