package storage

type Seeder interface {
	Seed(amount int) error
}
