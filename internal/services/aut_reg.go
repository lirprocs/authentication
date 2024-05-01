package services

import (
	"aut_reg/internal/jwt"
	"aut_reg/internal/models"
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Auth struct {
	usrSaver  UserSaver
	usrGetter UserGetter
	tokenTTL  time.Duration
}

type UserSaver interface {
	SaveUser(ctx context.Context, email, username string, password []byte) (int64, error)
}

type UserGetter interface {
	GetUser(ctx context.Context, username string) (models.User, error)
	//IsAdmin(ctx context.Context, userID int64) (bool, error)
}

func New(usrSaver UserSaver, usrProvider UserGetter, tokenTTL time.Duration) *Auth {
	return &Auth{
		usrSaver:  usrSaver,
		usrGetter: usrProvider,
		tokenTTL:  tokenTTL}
}

func (a *Auth) RegisterUser(ctx context.Context, email, username, password string) (int64, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("Error: %w", err)
	}
	id, err := a.usrSaver.SaveUser(ctx, email, username, hashPassword)
	if err != nil {
		return 0, fmt.Errorf("Error: %w", err)
	}

	return id, nil
}

func (a *Auth) LoginUser(ctx context.Context, username, password string) (string, error) {
	user, err := a.usrGetter.GetUser(ctx, username)
	if err != nil {
		return "", fmt.Errorf("Error: %w", err)
	}

	if err = bcrypt.CompareHashAndPassword(user.Password, []byte(password)); err != nil {
		return "", fmt.Errorf("Error: %w", err)
	}
	token, err := jwt.NewToken(user, a.tokenTTL)
	if err != nil {
		return "", fmt.Errorf("Error: %w", err)
	}

	return token, nil
}
