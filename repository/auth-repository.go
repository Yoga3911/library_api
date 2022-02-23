package repository

import (
	"context"
	"fmt"
	"project_restapi/models"
	"project_restapi/sql"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type AuthR interface {
	InsertData(ctx context.Context, user models.Register, hash string) error
	VerifyData(ctx context.Context, email string) (models.User, error)
	CheckDuplicate(ctx context.Context, name string, email string) error
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

func (a *authR) VerifyData(ctx context.Context, email string) (models.User, error) {
	var usr models.User
	var (
		createT, updateT time.Time
	)

	err := a.db.QueryRow(ctx, sql.VerifyCredential, email).Scan(&usr.ID, &usr.Name, &usr.Email, &usr.Password, &usr.GenderID, &usr.RoleID, &usr.Coin, &usr.IsActive, &createT, &updateT, &usr.Image)
	if err != nil {
		return usr, err
	}

	usr.CreateAt = createT.Format("02-01-2006")
	usr.UpdateAt = updateT.Format("02-01-2006")

	return usr, nil
}

func (a *authR) CheckDuplicate(ctx context.Context, name string, email string) error {
	var (
		nameC, emailC int8
	)
	
	err := a.db.QueryRow(ctx, sql.RegisterVal, name, email).Scan(&nameC, &emailC)
	if err != nil {
		return err
	}

	if nameC != 0 {
		return fmt.Errorf("duplicate name")
	}

	if emailC != 0 {
		return fmt.Errorf("duplicate email")
	}

	return nil
}

func (a *authR) UpdateActive(ctx context.Context, email string) error {
	_, err := a.db.Exec(ctx, sql.UpdateActive, email)
	
	return err
}