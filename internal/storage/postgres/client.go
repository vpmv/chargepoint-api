package postgres

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/vpmv/chargepoint-api/internal/dto"
	env "github.com/vpmv/chargepoint-api/pkg/dotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func NewClient(config Config, logger *logrus.Logger) (*Client, error) {
	var (
		retries int
		err     error
		db      *gorm.DB
	)

	for retries < 5 {
		db, err = gorm.Open(postgres.Open(config.DSN()))
		if err == nil {
			break
		}

		logger.Info(`Waiting to retry PostgreSQL connect...`)
		time.Sleep(5 * time.Second)
		retries++
	}

	if err != nil {
		return nil, err
	}

	return &Client{
		db:  db,
		log: logger,
	}, err
}

type Client struct {
	db   *gorm.DB
	log  *logrus.Logger
	seed bool
}

func (c Client) GetChargePoints(page, pageSize int) (dto.ChargePoints, error) {
	var chargePoints []*ChargePoint

	tx := c.paginate(c.db, page, pageSize)
	tx.Find(&chargePoints)

	if err := tx.Error; err != nil {
		return nil, err
	}

	return chargePointsToDTOSlice(chargePoints), nil
}

func (c Client) GetChargePoint(vendorID string) (*dto.ChargePoint, error) {
	var cp ChargePoint
	if err := c.db.First(&cp, `vendor_id = ?`, vendorID).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return nil, err
		} else {
			return nil, nil
		}
	}

	return chargePointToDTO(&cp), nil
}

func (c Client) CreateChargePoints(chargePoints dto.ChargePoints) (dto.ChargePoints, error) {
	models := dtoSliceToModel(chargePoints)
	err := c.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "vendor_id"},
		},
		DoUpdates: clause.AssignmentColumns([]string{
			"name",
			"point",
			"status",
			"updated_at",
		}),
	}).Create(&models).Error

	return chargePointsToDTOSlice(models), err
}

func (c Client) DeleteChargePoint(vendorID string) error {
	return c.db.Delete(&ChargePoint{}, "vendor_id = ?", vendorID).Error
}

func (c Client) GetChargePointByLocation(latitude, longitude float64, radiusKm int) (dto.ChargePoints, error) {
	var chargePoints []*ChargePoint

	if err := c.db.Raw(`
        SELECT
			id,
			name,
			status,
			vendor_id,
			point,
            ST_Distance(
                point,
                ST_SetSRID(ST_Point(@lon, @lat), 4326)::geography
            ) AS distance
        FROM charge_points
        WHERE 
            ST_DWithin(
				point,
				ST_SetSRID(ST_Point(@lon, @lat), 4326)::geography,
				@radiusMeter
        	) 
          AND deleted_at IS NULL
        ORDER BY distance
    `,
		map[string]interface{}{
			`lon`:         longitude,
			`lat`:         latitude,
			`radiusMeter`: radiusKm * 1000,
		},
	).Scan(&chargePoints).Error; err != nil {
		return nil, err
	}

	return chargePointsToDTOSlice(chargePoints), nil
}

func (c Client) Migrate() error {
	if err := c.db.AutoMigrate(&ChargePoint{}); err != nil {
		return err
	}

	if c.seed {
		seeder := NewSeeder(c, c.log)
		return seeder.Seed(env.GetInt(`SEED_COUNT`, 100))
	}

	return nil
}

func (c *Client) MustSeed() {
	c.seed = true
}

func (c Client) paginate(tx *gorm.DB, page, limit int) *gorm.DB {
	offset := limit * (page - 1)
	return tx.Offset(offset).Limit(limit)
}
