package cp

import (
	"testing"
	
	"github.com/stretchr/testify/assert"
	
	"github.com/creamsensation/cp/internal/tests"
	"github.com/creamsensation/quirk"
)

func TestUser(t *testing.T) {
	db := tests.CreateDatabaseConnection(t)
	if db == nil {
		assert.Fail(t, "connection doesn't exist")
		return
	}
	CreateUsersTable(quirk.New(db))
	um := CreateUserManager(db, 0, "")
	t.Cleanup(
		func() {
			DropUsersTable(quirk.New(db))
		},
	)
	t.Run(
		"create", func(t *testing.T) {
			id := um.Create(
				User{
					Active:   true,
					Roles:    []string{"owner"},
					Email:    "dominik@linduska.dev",
					Password: "123456789",
				},
			)
			assert.True(t, id > 0)
		},
	)
	t.Run(
		"get", func(t *testing.T) {
			var r User
			um.Get(&r)
			assert.True(t, r.Id > 0)
		},
	)
	t.Run(
		"update", func(t *testing.T) {
			data := User{
				Roles: []string{"admin"},
			}
			data.WithColumns("roles")
			um.Update(data)
			var r User
			um.Get(&r)
			assert.Equal(t, []string{"admin"}, r.Roles)
		},
	)
	t.Run(
		"disable enable", func(t *testing.T) {
			var r User
			um.Get(&r)
			assert.True(t, r.Active)
			um.Disable()
			um.Get(&r)
			assert.False(t, r.Active)
			um.Enable()
			um.Get(&r)
			assert.True(t, r.Active)
		},
	)
}
