package repository

import (
	"context"
	"project_restapi/models"
	"project_restapi/sql"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type AuthR interface {
	InsertData(ctx context.Context, user models.Register, hash string) error
	VerifyData(ctx context.Context, email string) pgx.Row
	CheckDuplicate(ctx context.Context, name string, email string) pgx.Row
	UpdateActive(ctx context.Context, email string) error
}

type authR struct {
	db *pgxpool.Pool
}

func NewAuthR(db *pgxpool.Pool) AuthR {
	return &authR{db: db}
}

func (a *authR) InsertData(ctx context.Context, user models.Register, hash string) error {
	_, err := a.db.Exec(ctx, sql.CreateUser, user.Name, user.Email, hash, user.GenderID)

	return err
}

func (a *authR) VerifyData(ctx context.Context, email string) pgx.Row {
	pg := a.db.QueryRow(ctx, sql.VerifyCredential, email)

	return pg
}

func (a *authR) CheckDuplicate(ctx context.Context, name string, email string) pgx.Row {
	pg := a.db.QueryRow(ctx, sql.RegisterVal, name, email)

	return pg
}

func (a *authR) UpdateActive(ctx context.Context, email string) error {
	_, err := a.db.Exec(ctx, sql.UpdateActive, email)
	
	return err
}