package postgres

import (
	"github.com/sirupsen/logrus"
	"github.com/vpmv/chargepoint-api/internal/dto"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (c Client) GetChargePoints(page, pageSize int) (*dto.ChargePoints, error) {
	var chargePoints []*ChargePoint
	tx := c.db.Find(&chargePoints)
	c.paginate(tx, page, pageSize)

	if err := tx.Error; err != nil {
		return nil, err
	}

	cps := dto.ChargePoints(chargePointsToDTOSlice(chargePoints))
	return &cps, nil
}

func (c Client) GetChargePoint(id string) (*dto.ChargePoint, error) {
	panic("implement me")
}

func (c Client) CreateChargePoints(chargePoints *dto.ChargePoints) (*dto.ChargePoints, error) {
	models := dtoSliceToModel(*chargePoints)
	err := c.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "vendor_id"},
		},
		DoUpdates: clause.AssignmentColumns([]string{
			"name",
			"latitude",
			"longitude",
			"status",
			"updated_at",
		}),
	}).Create(&models).Error

	cps := dto.ChargePoints(chargePointsToDTOSlice(models))
	return &cps, err
}

func (c Client) DeleteChargePoint(id string) error {
	panic("implement me")
}

func (c Client) GetChargePointByLocation(latitude, longitude float64, radiusKm int) *dto.ChargePoints {
	panic("implement me")
}

func (c Client) Migrate() error {
	return c.db.AutoMigrate(&ChargePoint{})
}

func (c Client) paginate(tx *gorm.DB, page, size int) {
	offset := size * (page - 1)
	if page == 1 {
		offset = 1
	}
	limit := size

	tx.Offset(offset).Limit(limit)
}
