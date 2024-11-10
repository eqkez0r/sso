package auth

import (
	"context"
	"errors"
	"github.com/eqkez0r/sso/internal/domain/models"
	"github.com/eqkez0r/sso/internal/lib/jwt"
	"github.com/eqkez0r/sso/internal/logger"
	"github.com/eqkez0r/sso/internal/services/storage"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserSaver interface {
	SaveUser(ctx context.Context, email string, passhash []byte) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAppID       = errors.New("invalid app id")
	ErrUserExists         = errors.New("user already exists")
)

type Auth struct {
	logger       logger.Logger
	userProvider UserProvider
	userSaver    UserSaver
	appProvider  AppProvider
	tokenTTL     time.Duration
}

func New(
	log logger.Logger,
	saver UserSaver,
	provider UserProvider,
	appProvider AppProvider,
	tokenTTL time.Duration) *Auth {
	return &Auth{
		logger:       log,
		userProvider: provider,
		userSaver:    saver,
		appProvider:  appProvider,
		tokenTTL:     tokenTTL,
	}
}

func (a Auth) Login(ctx context.Context, email, password string, appid int) (string, error) {
	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.logger.Warnf("%v", err)
			return "", ErrInvalidCredentials
		}
		return "", err
	}

	if err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.logger.Warnf("%v", err)
		return "", ErrInvalidCredentials
	}

	app, err := a.appProvider.App(ctx, appid)
	if err != nil {
		return "", err
	}

	a.logger.Info("Logged in")

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		a.logger.Warnf("%v", err)
		return "", err
	}
	return token, nil
}

func (a Auth) RegisterNewUser(ctx context.Context, email, password string) (userID int64, err error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		a.logger.Error(err)
		return -1, err
	}

	id, err := a.userSaver.SaveUser(ctx, email, hash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			a.logger.Warnf("%v", err)
			return -1, ErrUserExists
		}
		a.logger.Error(err)
		return -1, err
	}
	a.logger.Info("user registered")
	return id, nil
}

func (a Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	isAdmin, err := a.userProvider.IsAdmin(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			a.logger.Warnf("%v", err)
			return false, ErrInvalidAppID
		}
		a.logger.Error(err)
		return false, err
	}

	return isAdmin, nil
}
