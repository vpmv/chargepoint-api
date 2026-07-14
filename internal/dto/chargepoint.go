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

// type safety checks
var _ fuego.InTransformer = (*ChargePoint)(nil)
var _ fuego.OutTransformer = (*ChargePoint)(nil)
