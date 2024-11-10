package postgres

import (
	"context"
	"github.com/eqkez0r/sso/internal/domain/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgres struct {
	conn *pgxpool.Pool
}

func New(url string) (*postgres, error) {
	conn, err := pgxpool.New(context.Background(), url)
	if err != nil {
		return nil, err
	}
	if err = conn.Ping(context.Background()); err != nil {
		return nil, err
	}
	return &postgres{conn: conn}, nil
}

func (p *postgres) SaveUser(ctx context.Context, email string, passhash []byte) (int64, error) {
	var id int64
	err := p.conn.QueryRow(ctx, `INSERT INTO users(email, pass_hash) VALUES ($1, $2) RETURNING id`, email, passhash).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (p *postgres) User(ctx context.Context, email string) (models.User, error) {
	var user models.User
	if err := p.conn.QueryRow(ctx, `SELECT id, email, pass_hash FROM users WHERE email = $1`, email).
		Scan(&user.ID, &user.Email, &user.PassHash); err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (p *postgres) IsAdmin(ctx context.Context, userid int64) (bool, error) {
	var isAdmin bool
	if err := p.conn.QueryRow(ctx, `SELECT is_admin FROM users WHERE id = $1`, userid).Scan(&isAdmin); err != nil {
		return false, err
	}
	return isAdmin, nil
}

func (p *postgres) App(ctx context.Context, appid int) (models.App, error) {
	var app models.App
	if err := p.conn.QueryRow(ctx, `SELECT id, name, secret FROM app WHERE id = $1`, appid).Scan(&app.ID, &app.Name, &app.Secret); err != nil {
		return models.App{}, err
	}
	return app, nil
}

func (p *postgres) Close() {}
