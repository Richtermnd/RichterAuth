package auth

import (
	"context"
	"log/slog"
	"time"

	"github.com/Richtermnd/RichterAuth/internal/domain/models"
	"github.com/Richtermnd/RichterAuth/internal/domain/requests"
	"github.com/Richtermnd/RichterAuth/internal/domain/responses"
	"github.com/Richtermnd/RichterAuth/internal/errs"
)

type UserRepo interface {
	SaveUser(ctx context.Context, user models.User) (int, error)
	SaveConfirmKey(ctx context.Context, confirmKey models.ConfirmKey, userId int) error
	UserByEmail(ctx context.Context, email string) (models.User, error)
	UserById(ctx context.Context, id int) (models.User, error)
	ConfirmUser(ctx context.Context, id int) error
}

type Service struct {
	log      *slog.Logger
	userRepo UserRepo
}

func New(log *slog.Logger, userRepo UserRepo) *Service {
	return &Service{log: log, userRepo: userRepo}
}

func (s *Service) Register(ctx context.Context, registerRequest requests.Register) error {
	const op = "auth.service.RegisterUser"
	log := s.log.With(slog.String("op", op))

	// check exists
	user, err := s.userRepo.UserByEmail(ctx, registerRequest.Email)
	if err == nil {
		log.Error("user already exists", "err", err)
		return errs.ErrBadRequest("User already exists", nil)
	}
	if (user != models.User{}) {
		log.Error("user already exists", "err", err)
		return errs.ErrBadRequest("User already exists", nil)
	}

	// check passwords
	if registerRequest.Password != registerRequest.RepeatPassword {
		log.Error("passwords do not match", err)
		return errs.ErrBadRequest("Passwords do not match", nil)
	}

	// hash password
	hashedPassword, err := hashPassword(registerRequest.Password)
	if err != nil {
		log.Error("failed to hash password", err)
		return errs.ErrInternal(err)
	}

	// create user
	user = models.User{
		Username:       registerRequest.Username,
		Email:          registerRequest.Email,
		HashedPassword: hashedPassword,
		IsActive:       false,
		Role:           "user",
	}
	id, err := s.userRepo.SaveUser(ctx, user)
	if err != nil {
		log.Error("failed to save user", err)
		return errs.ErrInternal(err)
	}

	// generate confirm key
	confirmKey := models.ConfirmKey{
		UserId:     id,
		ConfirmKey: generateConfirmKey(10),
		ExpiredAt:  time.Now().Add(10 * time.Hour),
	}

	// save confirm key
	err = s.userRepo.SaveConfirmKey(ctx, confirmKey, id)
	if err != nil {
		log.Error("failed to save confirm key", "err", err)
		return errs.ErrInternal(err)
	}
	return nil
}

func (s *Service) Login(ctx context.Context, loginRequest requests.Login) (responses.Token, error) {
	const op = "auth.service.Login"
	log := s.log.With(slog.String("op", op))

	log.Info("Login user")
	user, err := s.userRepo.UserByEmail(ctx, loginRequest.Email)
	if err != nil {
		log.Error("failed to get user", "err", err)
		return responses.Token{}, err
	}

	// check password
	if !verifyPassword(user.HashedPassword, loginRequest.Password) {
		log.Error("invalid password")
		return responses.Token{}, errs.ErrBadRequest("Invalid password", nil)
	}

	// generate token
	token := generateToken(user.Id, user.Email, user.Role)
	return responses.Token{Token: token}, nil

}

func (s *Service) Confirm(ctx context.Context, confirmRequest requests.Confirm) error {
	const op = "auth.service.Confirm"
	log := s.log.With(slog.String("op", op))

	user, err := s.userRepo.UserById(ctx, confirmRequest.Id)
	if err != nil {
		log.Error("failed to get user", "err", err)
		return err
	}
	if user.ConfirmKey.ExpiredAt.Before(time.Now()) {
		log.Error("confirm key expired", "err", err)
		return errs.ErrBadRequest("Confirm key expired", nil)
	}

	if user.ConfirmKey.ConfirmKey != confirmRequest.Key {
		log.Error("invalid confirm key", "err", err)
		return errs.ErrBadRequest("Invalid confirm key", nil)
	}

	err = s.userRepo.ConfirmUser(ctx, confirmRequest.Id)
	if err != nil {
		log.Error("failed to confirm user", "err", err)
		return err
	}
	return nil
}
