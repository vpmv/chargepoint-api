package postgres

import (
	"github.com/sirupsen/logrus"
	"github.com/vpmv/chargepoint-api/internal/dto"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewClient(config Config, logger *logrus.Logger) (*Client, error) {
	if db, err := gorm.Open(postgres.Open(config.DSN()), &gorm.Config{}); err == nil {
		return &Client{db, logger}, err
	} else {
		return nil, err
	}
}

type Client struct {
	db  *gorm.DB
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
	return c.db.AutoMigrate(&ChargePoint{})
}
