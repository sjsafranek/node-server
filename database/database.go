package database

import (
	"errors"

	"github.com/sjsafranek/node-server/config"
	"github.com/sjsafranek/node-server/database/boltdb"
	"github.com/sjsafranek/node-server/models"
)

func New(options config.DatabaseConfiguration) (models.Database, error) {
	if "" == options.Type {
		options.Type = "disk"
	}

	connectionString, err := options.GetDatabaseConnectionString()
	if nil != err {
		return nil, err
	}

	switch options.Type {
	case "bolt":
		return boltdb.New(connectionString)
	default:
		return nil, errors.New("Unknown database type")
	}
}
