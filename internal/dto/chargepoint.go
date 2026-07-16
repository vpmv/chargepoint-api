package dto

import (
	"context"
	"strings"

	"github.com/go-fuego/fuego"
)

type Location struct {
	Latitude  float64 `json:"latitude" validate:"required,min=-90,max=90"`
	Longitude float64 `json:"longitude" validate:"required,min=-180,max=180"`
}

type ChargePoint struct {
	VendorId string   `json:"id" validate:"required"`
	Name     string   `json:"name" validate:"required"`
	Location Location `json:"location" validate:"required"`
	Status   string   `json:"status" validate:"oneOf=OFFLINE AVAILABLE OCCUPIED"`
	Distance float64  `json:"distance,omitempty"` // output only
}

func (c *ChargePoint) InTransform(ctx context.Context) error {
	c.Status = strings.ToLower(c.Status)
	return nil
}

func (c *ChargePoint) OutTransform(ctx context.Context) error {
	c.Status = strings.ToUpper(c.Status)
	return nil
}

type ChargePoints []*ChargePoint

func (s ChargePoints) InTransform(ctx context.Context) error {
	for i := range s {
		_ = s[i].InTransform(ctx)
	}
	return nil
}

func (s ChargePoints) OutTransform(ctx context.Context) error {
	for i := range s {
		_ = s[i].OutTransform(ctx)
	}
	return nil
}

// type safety checks
var _ fuego.InTransformer = (*ChargePoint)(nil)
var _ fuego.OutTransformer = (*ChargePoint)(nil)

var _ fuego.InTransformer = (*ChargePoints)(nil)
var _ fuego.OutTransformer = (*ChargePoints)(nil)
