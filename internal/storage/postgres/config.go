package postgres

import "fmt"

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DB       string
}

func (cfg Config) DSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		cfg.Host, cfg.User, cfg.Password, cfg.DB, cfg.Port,
	)
}
