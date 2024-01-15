package tests

import (
	"testing"
	
	"github.com/stretchr/testify/assert"
	
	"github.com/creamsensation/cp/internal/config"
	"github.com/creamsensation/quirk"
)

var Database = config.Database{
	Driver:   "postgres",
	Host:     "localhost",
	Port:     5432,
	User:     "cream",
	Password: "cream",
	Dbname:   "test",
	Ssl:      "disable",
}

func CreateDatabaseConnection(t *testing.T) *quirk.DB {
	db, err := quirk.Connect(
		quirk.WithPostgres(),
		quirk.WithHost(Database.Host),
		quirk.WithPort(Database.Port),
		quirk.WithDbname(Database.Dbname),
		quirk.WithUser(Database.User),
		quirk.WithPassword(Database.Password),
		quirk.WithSslDisable(),
	)
	assert.NoError(t, err)
	assert.NoError(t, db.Ping())
	return db
}
