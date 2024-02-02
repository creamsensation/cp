package cp

import (
	"fmt"

	"github.com/creamsensation/cp/internal/config"
	"github.com/creamsensation/cp/internal/connect"
	"github.com/creamsensation/quirk"
)

var (
	pgUserFields = []quirk.Field{
		{Name: quirk.Id, Props: "serial primary key"},
		{Name: UserActive, Props: "bool not null default false"},
		{Name: UserRoles, Props: "varchar[]"},
		{Name: UserEmail, Props: "varchar(255) not null"},
		{Name: UserPassword, Props: "varchar(128) not null"},
		{Name: UserTfa, Props: "bool not null default false"},
		{Name: UserTfaSecret, Props: "varchar(255)"},
		{Name: UserTfaCodes, Props: "varchar(255)"},
		{Name: UserTfaUrl, Props: "varchar(255)"},
		{Name: quirk.Vectors, Props: "tsvector not null default ''"},
		{Name: UserLastActivity, Props: "timestamp not null default current_timestamp"},
		{Name: quirk.CreatedAt, Props: "timestamp not null default current_timestamp"},
		{Name: quirk.UpdatedAt, Props: "timestamp not null default current_timestamp"},
	}
)

func CreateMigrationsConnections(config config.Databases) map[string]*quirk.DB {
	result := make(map[string]*quirk.DB)
	if len(config) == 0 {
		return result
	}
	for name, dbConfig := range config {
		result[name] = connect.Database(dbConfig)
	}
	return result
}

func CreateUsersTable(q *quirk.Quirk, customFields ...quirk.Field) {
	fields := make([]quirk.Field, 0)
	switch q.DB.DriverName() {
	case quirk.Postgres:
		for _, f := range pgUserFields {
			fields = append(fields, f)
		}
	}
	fields = quirk.MergeFields(fields, customFields)
	q.Q(
		fmt.Sprintf(
			`CREATE TABLE IF NOT EXISTS %s (%s)`,
			usersTable,
			quirk.CreateTableStructure(fields),
		),
	)
	q.MustExec()
}

func DropUsersTable(q *quirk.Quirk) {
	q.Q(fmt.Sprintf(`DROP TABLE IF EXISTS %s CASCADE`, usersTable)).MustExec()
}
