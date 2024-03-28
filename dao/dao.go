//go:build !solution

package dao

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"
)

func CreateDao(ctx context.Context, connString string) (Dao, error) {
	source, err := pgx.Connect(ctx, connString)

	if err != nil {
		return nil, err
	}

	_, err = source.Exec(ctx, `
   CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
   )
  `)

	if err != nil {
		return nil, err
	}

	dao := &PgxDao{conn: source}

	return dao, nil
}

func (d *PgxDao) Create(ctx context.Context, u *User) (UserID, error) {
	var id UserID

	err := d.conn.QueryRow(ctx, `
   INSERT INTO users (name) VALUES ($1) RETURNING id
  `, u.Name).Scan(&id)

	if err != nil {
		return 0, err
	}

	u.ID = id

	return id, nil
}

func (d *PgxDao) Update(ctx context.Context, u *User) error {
	exec, err := d.conn.Exec(ctx, `
   UPDATE users SET name=$1 WHERE id=$2
  `, u.Name, u.ID)

	if err != nil {
		return err
	}

	if exec.RowsAffected() == 0 {
		return fmt.Errorf("can't find user")
	}

	return nil
}

func (d *PgxDao) Delete(ctx context.Context, id UserID) error {
	_, err := d.conn.Exec(ctx, `
   DELETE FROM users WHERE id=$1
  `, id)

	return err
}

func (d *PgxDao) Lookup(ctx context.Context, id UserID) (User, error) {
	var name string
	sample := User{}

	err := d.conn.QueryRow(ctx, `
   SELECT name FROM users WHERE id=$1
  `, id).Scan(&name)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return sample, fmt.Errorf("can't find user: %w", sql.ErrNoRows)
		}

		return sample, err
	}

	return User{ID: id, Name: name}, nil
}

func (d *PgxDao) List(ctx context.Context) ([]User, error) {
	rows, err := d.conn.Query(ctx, "SELECT id, name FROM users")

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	users := make([]User, 0)

	for rows.Next() {
		var user User
		e := rows.Scan(&user.ID, &user.Name)

		if e != nil {
			return nil, e
		}

		users = append(users, user)
	}

	return users, rows.Err()
}

func (d *PgxDao) Close() error {
	_ = d.conn.Close(context.Background())

	return nil
}

type PgxDao struct {
	conn *pgx.Conn
}
