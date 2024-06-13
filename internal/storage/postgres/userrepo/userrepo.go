package userrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/Richtermnd/RichterAuth/internal/domain/models"
	"github.com/Richtermnd/RichterAuth/internal/errs"
	"github.com/jmoiron/sqlx"
)

var roles = make(map[string]int)

type UserRepo struct {
	db      *sqlx.DB
	builder sq.StatementBuilderType
}

func New(db *sqlx.DB) *UserRepo {
	repo := &UserRepo{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
	repo.uploadRoles(context.Background())
	return repo
}

func (r *UserRepo) SaveConfirmKey(ctx context.Context, confirmKey models.ConfirmKey, userId int) error {
	stmt, args, err := r.builder.Insert("confirm_keys").
		Columns("user_id", "confirm_key", "expired_at").
		Values(userId, confirmKey.ConfirmKey, confirmKey.ExpiredAt).
		ToSql()
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return errs.ErrInternal(err)
	}

	_, err = r.db.ExecContext(ctx, stmt, args...)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return errs.ErrInternal(err)
	}
	return nil
}

func (r *UserRepo) SaveUser(ctx context.Context, user models.User) (int, error) {
	stmt, args, err := r.builder.Insert("users").
		Columns("username", "email", "hashed_password", "role_id").
		Values(user.Username, user.Email, user.HashedPassword, 1).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return 0, errs.ErrInternal(err)
	}
	var id int
	err = r.db.GetContext(ctx, &id, stmt, args...)
	if err != nil {
		return 0, errs.ErrInternal(err)
	}
	return id, nil
}

func (r *UserRepo) User(ctx context.Context, email string) (models.User, error) {
	// Get user
	stmt, args, err := r.builder.
		Select(
			"u.id as id",
			"u.username as username",
			"u.email as email",
			"u.hashed_password as hashed_password",
			"u.registered_at as registered_at",
			"u.is_active as is_active",
			"u.is_confirmed as is_confirmed",
			"r.name as role",
		).
		From("users as u").
		Join("roles as r on r.id = u.role_id").
		Where(sq.Eq{"u.email": email}).
		ToSql()
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return models.User{}, errs.ErrInternal(err)
	}

	var user models.User
	err = r.db.GetContext(ctx, &user, stmt, args...)
	// Handle not found
	if errors.Is(err, sql.ErrNoRows) {
		return models.User{}, errs.ErrNotFound("user not found", err)
	} else if err != nil {
		fmt.Printf("err: %v\n", err)
		return models.User{}, errs.ErrInternal(err)
	}

	// Get confirm key
	stmt, args, err = r.builder.
		Select(
			"confirm_key",
			"expired_at",
		).
		From("confirm_keys").
		Where(sq.Eq{"user_id": user.Id}).
		ToSql()
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return models.User{}, errs.ErrInternal(err)
	}
	var confirmKey models.ConfirmKey
	err = r.db.GetContext(ctx, &confirmKey, stmt, args...)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return models.User{}, errs.ErrInternal(err)
	}
	user.ConfirmKey = confirmKey
	return user, nil
}

func (r *UserRepo) uploadRoles(ctx context.Context) {
	rows, err := r.db.QueryContext(ctx, "SELECT id, name FROM roles")
	if err != nil {
		panic("could not upload roles")
	}
	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			panic("could not upload roles")
		}
		roles[name] = id
	}
}
