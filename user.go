package cp

import (
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/matthewhartstonge/argon2"

	"github.com/creamsensation/quirk"
)

type UserManager interface {
	Get(id ...int) User
	Create(r User) int
	Update(r User, columns ...string)
	UpdatePassword(actualPassword, newPassword string) error
	ForceUpdatePassword(newPassword string) error
	Enable(id ...int)
	Disable(id ...int)
}

type UserField struct {
	Name  string
	Value any
}

type userManager struct {
	db         *quirk.DB
	id         int
	email      string
	driverName string
	data       map[string]any
}

type User struct {
	Id           int            `json:"id"`
	Active       bool           `json:"active"`
	Roles        []string       `json:"roles"`
	Email        string         `json:"email"`
	Password     string         `json:"password"`
	Tfa          bool           `json:"tfa"`
	TfaSecret    sql.NullString `json:"tfaSecret"`
	TfaCodes     sql.NullString `json:"tfaCodes"`
	TfaUrl       sql.NullString `json:"tfaUrl"`
	LastActivity time.Time      `json:"lastActivity"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`

	columns []string
}

const (
	UserActive       = "active"
	UserRoles        = "roles"
	UserEmail        = "email"
	UserPassword     = "password"
	UserTfa          = "tfa"
	UserTfaSecret    = "tfa_secret"
	UserTfaCodes     = "tfa_codes"
	UserTfaUrl       = "tfa_url"
	UserLastActivity = "last_activity"
)

const (
	usersTable  = "users"
	paramPrefix = "@"
)

const (
	operationInsert = "insert"
	operationUpdate = "update"
)

var (
	ErrorMissingUser      = errors.New("user doesn't exist")
	ErrorMismatchPassword = errors.New("passwords aren't equal")
	ErrorHash             = errors.New("hash failed")
)

var (
	argon = argon2.DefaultConfig()
)

func CreateUserManager(db *quirk.DB, id int, email string) UserManager {
	return &userManager{
		db:         db,
		email:      email,
		id:         id,
		data:       make(map[string]any),
		driverName: db.DriverName(),
	}
}

func (u *userManager) Get(id ...int) User {
	if len(id) > 0 {
		u.id = id[0]
	}
	var r User
	if u.id == 0 && u.email == "" {
		return r
	}
	quirk.New(u.db).Q(`SELECT *`).
		Q(fmt.Sprintf(`FROM %s`, usersTable)).
		If(u.id > 0, `WHERE id = ?`, u.id).
		If(u.id == 0, `WHERE email = ?`, u.email).
		Q(`LIMIT 1`).
		MustExec(&r)
	clear(u.data)
	return r
}

func (u *userManager) Create(r User) int {
	if u.id != 0 {
		return u.id
	}
	u.readData(operationInsert, r, []string{})
	columns, placeholders := u.insertValues()
	quirk.New(u.db).Q(fmt.Sprintf(`INSERT INTO %s`, usersTable)).
		Q(fmt.Sprintf(`(%s)`, columns)).
		Q(fmt.Sprintf(`VALUES (%s)`, placeholders), u.args()...).
		Q(`RETURNING id`).
		MustExec(&u.id)
	u.email = r.Email
	clear(u.data)
	return u.id
}

func (u *userManager) Update(r User, columns ...string) {
	if u.id == 0 && u.email == "" {
		return
	}
	u.readData(operationUpdate, r, columns)
	quirk.New(u.db).Q(fmt.Sprintf(`UPDATE %s`, usersTable)).
		Q(fmt.Sprintf(`SET %s`, u.updateValues()), u.args()...).
		If(u.id > 0, `WHERE id = ?`, u.id).
		If(u.id == 0, `WHERE email = ?`, u.email).
		MustExec()
	clear(u.data)
}

func (u *userManager) UpdatePassword(actualPassword, newPassword string) error {
	if u.id == 0 && u.email == "" {
		return ErrorMissingUser
	}
	user := u.Get()
	if ok, err := argon2.VerifyEncoded([]byte(actualPassword), []byte(user.Password)); !ok || err != nil {
		return ErrorMismatchPassword
	}
	err := quirk.New(u.db).Q(fmt.Sprintf(`UPDATE %s`, usersTable)).
		Q(`SET password = ?`, u.hashPassword(newPassword)).
		If(u.id > 0, `WHERE id = ?`, u.id).
		If(u.id == 0, `WHERE email = ?`, u.email).
		Exec()
	clear(u.data)
	return err
}

func (u *userManager) ForceUpdatePassword(newPassword string) error {
	if u.id == 0 && u.email == "" {
		return ErrorMissingUser
	}
	err := quirk.New(u.db).Q(fmt.Sprintf(`UPDATE %s`, usersTable)).
		Q(`SET password = ?`, u.hashPassword(newPassword)).
		If(u.id > 0, `WHERE id = ?`, u.id).
		If(u.id == 0, `WHERE email = ?`, u.email).
		Exec()
	clear(u.data)
	return err
}

func (u *userManager) Enable(id ...int) {
	if len(id) > 0 {
		u.id = id[0]
	}
	if u.id == 0 && u.email == "" {
		return
	}
	quirk.New(u.db).Q(fmt.Sprintf(`UPDATE %s`, usersTable)).
		Q(`SET active = true`).
		If(u.id > 0, `WHERE id = ?`, u.id).
		If(u.id == 0, `WHERE email = ?`, u.email).
		MustExec()
	clear(u.data)
}

func (u *userManager) Disable(id ...int) {
	if len(id) > 0 {
		u.id = id[0]
	}
	if u.id == 0 && u.email == "" {
		return
	}
	quirk.New(u.db).Q(fmt.Sprintf(`UPDATE %s`, usersTable)).
		Q(`SET active = false`).
		If(u.id > 0, `WHERE id = ?`, u.id).
		If(u.id == 0, `WHERE email = ?`, u.email).
		MustExec()
	clear(u.data)
}

func (u *userManager) readData(operation string, data User, columns []string) {
	columnsExist := len(columns) > 0
	if operation == operationInsert && slices.Contains(columns, quirk.Id) {
		u.data[quirk.Id] = data.Id
	}
	if !columnsExist || slices.Contains(columns, UserActive) {
		u.data[UserActive] = data.Active
	}
	if !columnsExist || slices.Contains(columns, UserEmail) {
		u.data[UserEmail] = data.Email
	}
	if !columnsExist || slices.Contains(columns, UserPassword) {
		u.data[UserPassword] = u.hashPassword(data.Password)
	}
	if !columnsExist || slices.Contains(columns, UserRoles) {
		u.data[UserRoles] = data.Roles
	}
	if !columnsExist || slices.Contains(columns, UserTfa) {
		u.data[UserTfa] = data.Tfa
	}
	if !columnsExist || slices.Contains(columns, UserTfaSecret) {
		u.data[UserTfaSecret] = data.TfaSecret.String
	}
	if !columnsExist || slices.Contains(columns, UserTfaCodes) {
		u.data[UserTfaCodes] = data.TfaCodes.String
	}
}

func (u *userManager) hashPassword(password string) string {
	hash, err := argon.HashEncoded([]byte(password))
	if err != nil {
		return password
	}
	return string(hash)
}

func (u *userManager) insertValues() (string, string) {
	columns := []string{quirk.Id}
	placeholders := []string{quirk.Default}
	for name := range u.data {
		columns = append(columns, name)
		placeholders = append(placeholders, paramPrefix+name)
	}
	switch u.driverName {
	case quirk.Postgres:
		if len(u.data) > 0 {
			columns = append(columns, quirk.Vectors)
			placeholders = append(placeholders, paramPrefix+quirk.Vectors)
		}
	}
	columns = append(columns, UserLastActivity)
	placeholders = append(placeholders, quirk.CurrentTimestamp)

	columns = append(columns, quirk.CreatedAt)
	placeholders = append(placeholders, quirk.CurrentTimestamp)

	columns = append(columns, quirk.UpdatedAt)
	placeholders = append(placeholders, quirk.CurrentTimestamp)
	return strings.Join(columns, ","), strings.Join(placeholders, ",")
}

func (u *userManager) args() []any {
	if len(u.data) == 0 {
		return []any{}
	}
	result := u.data
	vectors := make([]any, 0)
	for name, v := range u.data {
		if name == UserPassword {
			continue
		}
		vectors = append(vectors, v)
	}
	switch u.driverName {
	case quirk.Postgres:
		if len(vectors) > 0 {
			result[quirk.Vectors] = quirk.CreateTsVectors(vectors...)
		}
	}
	return []any{result}
}

func (u *userManager) updateValues() string {
	result := make([]string, 0)
	for column := range u.data {
		if column == quirk.Id {
			continue
		}
		result = append(result, fmt.Sprintf("%s = %s%s", column, paramPrefix, column))
	}
	result = append(result, fmt.Sprintf("%s = %s", UserLastActivity, quirk.CurrentTimestamp))
	result = append(result, fmt.Sprintf("%s = %s", quirk.UpdatedAt, quirk.CurrentTimestamp))
	switch u.driverName {
	case quirk.Postgres:
		vectors := make([]any, 0)
		for column, v := range u.data {
			if column == quirk.Id {
				continue
			}
			vectors = append(vectors, v)
		}
		if len(vectors) > 0 {
			result = append(result, fmt.Sprintf("%s = %v", quirk.Vectors, quirk.CreateTsVectors(vectors...).Value))
		}
	}
	return strings.Join(result, ",")
}
