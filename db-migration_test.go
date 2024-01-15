package cp

import (
	"testing"
	
	"github.com/stretchr/testify/assert"
	
	"github.com/creamsensation/cp/internal/config"
	"github.com/creamsensation/cp/internal/constant/naming"
	"github.com/creamsensation/cp/internal/tests"
	"github.com/creamsensation/quirk"
)

func TestDbMigration(t *testing.T) {
	var db *quirk.DB
	t.Cleanup(
		func() {
			if db == nil {
				assert.Fail(t, "connection doesn't exist")
				return
			}
			assert.Nil(t, db.Close())
		},
	)
	t.Run(
		"create migrations connections", func(t *testing.T) {
			connections := CreateMigrationsConnections(
				config.Databases{
					naming.Main: tests.Database,
				},
			)
			assert.Nil(t, connections[naming.Main].Ping())
			db = connections[naming.Main]
		},
	)
	t.Run(
		"create users table", func(t *testing.T) {
			if db == nil {
				assert.Fail(t, "connection doesn't exist")
				return
			}
			CreateUsersTable(quirk.New(db))
			var count int
			assert.Nil(
				t,
				quirk.New(db).Q(`select count(TABLE_NAME) from information_schema.tables where TABLE_NAME = 'users'`).Exec(&count),
			)
			assert.True(t, count == 1)
		},
	)
	t.Run(
		"drop users table", func(t *testing.T) {
			if db == nil {
				assert.Fail(t, "connection doesn't exist")
				return
			}
			DropUsersTable(quirk.New(db))
			var count int
			assert.Nil(
				t,
				quirk.New(db).Q(`select count(TABLE_NAME) from information_schema.tables where TABLE_NAME = 'users'`).Exec(&count),
			)
			assert.True(t, count == 0)
		},
	)
}
