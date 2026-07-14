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
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DB,
	)
}
