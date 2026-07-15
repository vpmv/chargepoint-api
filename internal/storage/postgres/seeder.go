package postgres

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/vpmv/chargepoint-api/internal/dto"
)

var vendorsIDs = []string{`SH`, `FN`, `GU`, `OK`, `ES`}

type ChargePointSeeder struct {
	client Client
	log    *logrus.Logger
	rng    *rand.Rand
}

// Seed seeds charge points into the database
// Param amount - number of charge points to seed per vendor
func (s ChargePointSeeder) Seed(amount int) error {
	s.log.Infoln(`Seeding ChargePoints...`)

	var points []*dto.ChargePoint

	for _, vendor := range vendorsIDs {
		points = make([]*dto.ChargePoint, 0, amount*10)
		for station := 1; station <= 10; station++ {
			status := s.rng.Intn(4) - 1
			for outlet := 0; outlet < amount; outlet++ {
				points = append(points, &dto.ChargePoint{
					VendorId: fmt.Sprintf(`%s%d-%d`, vendor, station, outlet),
					Name:     fmt.Sprintf("Mock charger %s%d-%d", vendor, station, outlet),
					Location: s.randomLocation(),
					Status:   StatusNames[chargePointStatus(status)],
				})
			}
		}

		s.log.Infof("Seeding %d charge points for vendor %s", len(points), vendor)
		if _, err := s.client.CreateChargePoints(points); err != nil {
			return err
		}
	}

	s.log.Infoln(`Done!`)
	return nil
}

// randomLocation generates a random point roughly within the EU region
func (s ChargePointSeeder) randomLocation() dto.Location {
	coordRange := func(min, max float64) float64 {
		return min + s.rng.Float64()*(max-min)
	}

	return dto.Location{
		Longitude: coordRange(-10, 30),
		Latitude:  coordRange(40, 60),
	}
}

// NewSeeder creates a new ChargePointSeeder instance
func NewSeeder(client Client, log *logrus.Logger) *ChargePointSeeder {

	return &ChargePointSeeder{
		client: client,
		log:    log,
		rng:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}
