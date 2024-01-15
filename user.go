package cp

import (
	"database/sql"
	"fmt"
	"slices"
	"strings"
	"time"
	
	"github.com/matthewhartstonge/argon2"
	
	"github.com/creamsensation/quirk"
)

type UserManager interface {
	Get(w UserWriter)
	Create(r UserReader) int
	Update(r UserReader)
	Enable()
	Disable()
}

type UserReader interface {
	GetId() int
	GetActive() bool
	GetRoles() []string
	GetEmail() string
	GetPassword() string
	GetTfa() bool
	GetTfaSecret() string
	GetTfaCodes() string
	GetCustomFields() []UserField
	GetColumns() []string
}

type UserWriter interface {
	GetColumns() []string
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
	UserColumnActive       = "active"
	UserColumnRoles        = "roles"
	UserColumnEmail        = "email"
	UserColumnPassword     = "password"
	UserColumnTfa          = "tfa"
	UserColumnTfaSecret    = "tfa_secret"
	UserColumnTfaCodes     = "tfa_codes"
	UserColumnTfaUrl       = "tfa_url"
	UserColumnLastActivity = "last_activity"
)

const (
	usersTable = "users"
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

func (u *userManager) Get(w UserWriter) {
	columns := "*"
	if len(w.GetColumns()) > 0 {
		columns = strings.Join(w.GetColumns(), ",")
	}
	quirk.New(u.db).Q(fmt.Sprintf(`SELECT %s`, columns)).
		Q(fmt.Sprintf(`FROM %s`, usersTable)).
		If(u.id > 0, `WHERE id = ?`, u.id).
		If(u.id == 0, `WHERE email = ?`, u.email).
		Q(`LIMIT 1`).
		MustExec(w)
	clear(u.data)
}

func (u *userManager) Create(r UserReader) int {
	if u.id != 0 {
		var usr User
		u.Get(&usr)
		return usr.Id
	}
	u.readData(r)
	columns, placeholders := u.insertValues()
	quirk.New(u.db).Q(fmt.Sprintf(`INSERT INTO %s`, usersTable)).
		Q(fmt.Sprintf(`(%s)`, columns)).
		Q(fmt.Sprintf(`VALUES (%s)`, placeholders), u.args()...).
		Q(`RETURNING id`).
		MustExec(&u.id)
	u.email = r.GetEmail()
	clear(u.data)
	return u.id
}

func (u *userManager) Update(r UserReader) {
	u.readData(r)
	quirk.New(u.db).Q(fmt.Sprintf(`UPDATE %s`, usersTable)).
		Q(fmt.Sprintf(`SET %s`, u.updateValues()), u.args()...).
		If(u.id > 0, `WHERE id = ?`, u.id).
		If(u.id == 0, `WHERE email = ?`, u.email).
		MustExec()
	clear(u.data)
}

func (u *userManager) Enable() {
	quirk.New(u.db).Q(fmt.Sprintf(`UPDATE %s`, usersTable)).
		Q(`SET active = true`).
		If(u.id > 0, `WHERE id = ?`, u.id).
		If(u.id == 0, `WHERE email = ?`, u.email).
		MustExec()
	clear(u.data)
}

func (u *userManager) Disable() {
	quirk.New(u.db).Q(fmt.Sprintf(`UPDATE %s`, usersTable)).
		Q(`SET active = false`).
		If(u.id > 0, `WHERE id = ?`, u.id).
		If(u.id == 0, `WHERE email = ?`, u.email).
		MustExec()
	clear(u.data)
}

func (u *userManager) readData(data UserReader) {
	columns := data.GetColumns()
	id := data.GetId()
	active := data.GetActive()
	email := data.GetEmail()
	password := data.GetPassword()
	roles := data.GetRoles()
	tfa := data.GetTfa()
	tfaSecret := data.GetTfaSecret()
	tfaCodes := data.GetTfaCodes()
	if id > 0 {
		u.data[quirk.Id] = id
	}
	if active || slices.Contains(columns, UserColumnActive) {
		u.data[UserColumnActive] = active
	}
	if len(email) > 0 || slices.Contains(columns, UserColumnEmail) {
		u.data[UserColumnEmail] = email
	}
	if len(password) > 0 || slices.Contains(columns, UserColumnPassword) {
		u.data[UserColumnPassword] = u.hashPassword(password)
	}
	if len(roles) > 0 || slices.Contains(columns, UserColumnRoles) {
		u.data[UserColumnRoles] = roles
	}
	if tfa || slices.Contains(columns, UserColumnTfa) {
		u.data[UserColumnTfa] = tfa
	}
	if len(tfaSecret) > 0 || slices.Contains(columns, UserColumnTfaSecret) {
		u.data[UserColumnTfaSecret] = tfaSecret
	}
	if len(tfaCodes) > 0 || slices.Contains(columns, UserColumnTfaCodes) {
		u.data[UserColumnTfaCodes] = tfaCodes
	}
	for _, f := range data.GetCustomFields() {
		if f.Value == nil && !slices.Contains(columns, f.Name) {
			continue
		}
		u.data[f.Name] = f.Value
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
		placeholders = append(placeholders, ":"+name)
	}
	switch u.driverName {
	case quirk.Postgres:
		if len(u.data) > 0 {
			columns = append(columns, quirk.Vectors)
			placeholders = append(placeholders, ":"+quirk.Vectors)
		}
	}
	columns = append(columns, UserColumnLastActivity)
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
		if name == UserColumnPassword {
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
		result = append(result, fmt.Sprintf("%s = :%s", column, column))
	}
	result = append(result, fmt.Sprintf("%s = %s", UserColumnLastActivity, quirk.CurrentTimestamp))
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

func (u User) GetId() int {
	return u.Id
}

func (u User) GetActive() bool {
	return u.Active
}

func (u User) GetRoles() []string {
	return u.Roles
}

func (u User) GetEmail() string {
	return u.Email
}

func (u User) GetPassword() string {
	return u.Password
}

func (u User) GetTfa() bool {
	return u.Tfa
}

func (u User) GetTfaSecret() string {
	return u.TfaSecret.String
}

func (u User) GetTfaCodes() string {
	return u.TfaCodes.String
}

func (u User) GetCustomFields() []UserField {
	return []UserField{}
}

func (u User) GetColumns() []string {
	return u.columns
}

func (u User) WithColumns(columns ...string) User {
	u.columns = columns
	return u
}
