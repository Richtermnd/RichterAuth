package models

import "time"

type User struct {
	Id             int       `db:"id"`
	Username       string    `db:"username"`
	Email          string    `db:"email"`
	HashedPassword string    `db:"hashed_password"`
	RegisteredAt   time.Time `db:"registered_at"`
	IsActive       bool      `db:"is_active"`
	IsConfirmed    bool      `db:"is_confirmed"`
	Role           string    `db:"role"`
	ConfirmKey     ConfirmKey
}

type ConfirmKey struct {
	UserId     int       `db:"user_id"`
	ConfirmKey string    `db:"confirm_key"`
	ExpiredAt  time.Time `db:"expired_at"`
}
