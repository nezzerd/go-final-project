package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPostgresConnection_InvalidConfig(t *testing.T) {
	cfg := Config{
		Host:     "invalid-host",
		Port:     "5432",
		User:     "user",
		Password: "pass",
		DBName:   "db",
	}

	db, err := NewPostgresConnection(cfg)
	assert.Error(t, err)
	assert.Nil(t, db)
}
