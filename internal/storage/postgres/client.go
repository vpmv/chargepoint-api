package postgres

import (
	"github.com/sirupsen/logrus"
	"github.com/vpmv/chargepoint-api/internal/dto"
)

func NewClient(config Config, logger *logrus.Logger) (*Client, error) {
	return &Client{}, nil
}

type Client struct {
	log *logrus.Logger
}

func (c Client) GetChargePoints() ([]*dto.ChargePoint, error) {
	panic("implement me")
}

func (c Client) GetChargePoint(id string) (*dto.ChargePoint, error) {
	panic("implement me")
}

func (c Client) CreateChargePoint(chargePoint *dto.ChargePoint) (*dto.ChargePoint, error) {
	panic("implement me")
}

func (c Client) DeleteChargePoint(id string) error {
	panic("implement me")
}

func (c Client) GetChargePointByLocation(latitude, longitude float64, radiusKm int) []*dto.ChargePoint {
	panic("implement me")
}

func (c Client) Migrate() error {
	panic("implement me")
}
