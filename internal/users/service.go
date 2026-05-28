package users

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/nikhil-sharma-dotcom/my-ecommerce-api/internal/adapters/postgresql/sqlc"
	myauth "github.com/nikhil-sharma-dotcom/my-ecommerce-api/internal/auth"
	"github.com/nikhil-sharma-dotcom/my-ecommerce-api/internal/config"
)

type Service struct {
	queries *sqlc.Queries
	config  *config.Config
}

func NewService(queries *sqlc.Queries, cfg *config.Config) *Service {
	return &Service{queries: queries, config: cfg}
}

func (s *Service) Register(ctx context.Context, email, password, firstName, lastName string) (sqlc.User, error) {
	_, err := s.queries.GetUserByEmail(ctx, email)
	if err == nil {
		return sqlc.User{}, errors.New("user already exists")
	}

	hash, err := myauth.HashPassword(password)
	if err != nil {
		return sqlc.User{}, err
	}

	return s.queries.CreateUser(ctx, sqlc.CreateUserParams{
		Email:        email,
		PasswordHash: hash,
		FirstName:    firstName,
		LastName:     lastName,
		Role:         pgtype.Text{String: "customer", Valid: true},
	})
}

func (s *Service) Login(ctx context.Context, email, password string) (string, sqlc.User, error) {
	user, err := s.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return "", sqlc.User{}, errors.New("invalid credentials")
	}

	if !myauth.CheckPasswordHash(password, user.PasswordHash) {
		return "", sqlc.User{}, errors.New("invalid credentials")
	}

	token, err := myauth.GenerateToken(int(user.ID), user.Email, user.Role.String, s.config.JWTSecret)
	if err != nil {
		return "", sqlc.User{}, err
	}

	return token, user, nil
}

func (s *Service) GetByID(ctx context.Context, id int64) (sqlc.User, error) {
	return s.queries.GetUserByID(ctx, id)
}
