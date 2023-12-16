package cp

import (
	"context"
	"testing"
	
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	
	cfg "github.com/creamsensation/cp/internal/config"
	"github.com/creamsensation/gox"
	"github.com/creamsensation/quirk"
)

var (
	testsLangs = cfg.Languages{
		"cs": cfg.Language{Enabled: true, Default: true},
		"en": cfg.Language{Enabled: true},
	}
)

type testComponent struct {
	Component
}

func (c *testComponent) Name() string {
	return "test"
}

func (c *testComponent) Model() {
}

func (c *testComponent) Node() gox.Node {
	return gox.Text("test")
}

func createDatabaseConnection() (*quirk.DB, error) {
	return quirk.Connect(
		quirk.WithPostgres(),
		quirk.WithHost("localhost"),
		quirk.WithPort(5432),
		quirk.WithDbname("test"),
		quirk.WithUser("cream"),
		quirk.WithPassword("cream"),
		quirk.WithSslDisable(),
	)
}

func createRedisConnection() (*redis.Client, error) {
	client := redis.NewClient(
		&redis.Options{
			Addr: "localhost:6379",
			DB:   10,
		},
	)
	return client, client.Ping(context.Background()).Err()
}

func createTestConnections(t *testing.T) (*quirk.DB, *redis.Client) {
	db, err := createDatabaseConnection()
	assert.NoError(t, err)
	assert.NoError(t, db.Ping())
	cache, err := createRedisConnection()
	assert.NoError(t, err)
	assert.NoError(t, cache.Ping(context.Background()).Err())
	return db, cache
}
