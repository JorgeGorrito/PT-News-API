package config

import (
	"fmt"
	"time"
)

type DatabaseConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	Database        string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=UTC",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
	)
}

func DefaultDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:            "localhost",
		Port:            3306,
		User:            "root",
		Password:        "root",
		Database:        "pt_news_api",
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5 * time.Minute,
	}
}
