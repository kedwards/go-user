package models

import (
	"context"
	"database/sql"
	"time"
)

// DBModel is the type for database connection values
type DBModel struct {
	DB *sql.DB
}

// Models is the wrapper for all models
type Models struct {
	DB DBModel
}

// NewModels returns a model type with database connection pool
func NewModels(db *sql.DB) Models {
	return Models{
		DB: DBModel{DB: db},
	}
}

// User is the type for users
type User struct {
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Username  string    `json:"username"`
	Role			string    `json:"role"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// GetUser gets all users
func (m *DBModel) GetUsers() ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var users []*User

	query := `
		select
			email, password, first_name, last_name, role, username, created_at, updated_at
		from
			users
	`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var u User
		err := rows.Scan(
			&u.Email,
			&u.Password,
			&u.FirstName,
			&u.LastName,
			&u.Role,
			&u.Username,
			&u.CreatedAt,
			&u.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &u)
	}

	return users, nil
}

// CreateUser creates a new user
func (m *DBModel) CreateUser(user User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		insert into users (email, password, first_name, last_name, role, username, created_at, updated_at)
		values
		($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := m.DB.ExecContext(ctx, stmt,
		user.Email,
		user.Password,
		user.FirstName,
		user.LastName,
		user.Role,
		user.Username,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

// UpdateUser updates a user
func (m *DBModel) UpdateUser(user User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
	  update users
		set password = $2, first_name = $3, last_name = $4, role = $5, username = $6, updated_at = $7
		where email = $1
	`

	_, err := m.DB.ExecContext(ctx, stmt,
		user.Email,
		user.Password,
		user.FirstName,
		user.LastName,
		user.Role,
		user.Username,
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

// DeleteUser deletes a user
func (m *DBModel) DeleteUser(userEmail string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
	  delete from users
		where email = $1
	`

	_, err := m.DB.ExecContext(ctx, stmt, userEmail)

	if err != nil {
		return err
	}

	return nil
}
