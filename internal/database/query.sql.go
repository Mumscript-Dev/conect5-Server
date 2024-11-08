// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: query.sql

package database

import (
	"context"
	"database/sql"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (id, email, username, profileIndex)
VALUES (?, ?, ?, ?)
RETURNING id, email, username, profileindex
`

type CreateUserParams struct {
	ID           interface{}
	Email        string
	Username     sql.NullString
	Profileindex sql.NullString
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.ID,
		arg.Email,
		arg.Username,
		arg.Profileindex,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Username,
		&i.Profileindex,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, email, username, profileindex FROM users WHERE email = ?
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Username,
		&i.Profileindex,
	)
	return i, err
}

const listUsers = `-- name: ListUsers :many
SELECT id, email, username, profileindex FROM users
`

func (q *Queries) ListUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, listUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Email,
			&i.Username,
			&i.Profileindex,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateUser = `-- name: UpdateUser :exec
UPDATE users
SET email = ?, username = ?, profileIndex = ?
WHERE email = ?
`

type UpdateUserParams struct {
	Email        string
	Username     sql.NullString
	Profileindex sql.NullString
	Email_2      string
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) error {
	_, err := q.db.ExecContext(ctx, updateUser,
		arg.Email,
		arg.Username,
		arg.Profileindex,
		arg.Email_2,
	)
	return err
}
