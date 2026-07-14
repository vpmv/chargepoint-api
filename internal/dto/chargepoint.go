package dto

import (
	"context"
	"strings"

	"github.com/go-fuego/fuego"
)

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type ChargePoint struct {
	ID       string   `json:"id"`
	VendorId string   `json:"vendor_id"`
	Name     string   `json:"name"`
	Location Location `json:"location"`
	Status   string   `json:"status"`
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
