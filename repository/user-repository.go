package repository

import (
	"context"
	"project_restapi/models"
	"project_restapi/sql"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type UserR interface {
	GetAll() (pgx.Rows, error)
	GetOne(ctx context.Context, id float64) pgx.Row
	Update(ctx context.Context, update models.Update, id float64) error
	ChangePassword(ctx context.Context, id float64, password string) error
	Delete(ctx context.Context, id float64) error
	CheckEmail(ctx context.Context, email string) pgx.Row
	CountQuantity(ctx context.Context, book_id string) pgx.Row
	CheckBook(ctx context.Context, book_id string, id float64) pgx.Row
	GetCoin(ctx context.Context, id float64) pgx.Row
	UpdateCoinUser(ctx context.Context, id float64, coin int) error
	TakeBook(ctx context.Context, book_id string, id float64) error
	GetById(ctx context.Context, id float64) (pgx.Rows, error)
	ValidateDeleteBook(ctx context.Context, book_id string, id float64) pgx.Row
	DeleteBookById(ctx context.Context, book_id string, id float64) error
	UserOneBook(ctx context.Context, id float64, book_id string) pgx.Row
	ReqAdmin(ctx context.Context, id float64) error
	UpdateRequest(ctx context.Context, answer models.Request, admin_id float64) error
	TokenReqAdmin(ctx context.Context, id uint64) pgx.Row
	CheckReqAdmin(ctx context.Context, id float64) pgx.Row
	GetAllRequest(ctx context.Context) (pgx.Rows, error)
}

type userR struct {
	db *pgxpool.Pool
}

func NewUserR(db *pgxpool.Pool) UserR {
	return &userR{
		db: db,
	}
}

func (u *userR) GetAll() (pgx.Rows, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	pg, err := u.db.Query(ctx, sql.GetAll)

	return pg, err
}

func (u *userR) GetOne(ctx context.Context, id float64) pgx.Row {
	pg := u.db.QueryRow(ctx, sql.GetOne, id)

	return pg
}

func (u *userR) Update(ctx context.Context, update models.Update, id float64) error {
	_, err := u.db.Exec(ctx, sql.UpdateUser, id, update.Name, update.Email, update.GenderID, update.Image)

	return err
}

func (u *userR) ChangePassword(ctx context.Context, id float64, password string) error {
	_, err := u.db.Exec(ctx, sql.ChangePass, id, password)
	
	return err
}

func (u *userR) Delete(ctx context.Context, id float64) error {
	_, err := u.db.Exec(ctx, sql.DeleteUser, id)

	return err
}

func (u *userR) CheckEmail(ctx context.Context, email string) pgx.Row {
	pg := u.db.QueryRow(ctx, sql.GetByEmail, email)

	return pg
}

func (u *userR) CountQuantity(ctx context.Context, book_id string) pgx.Row {
	pg := u.db.QueryRow(ctx, sql.CheckQuantity, book_id)
	
	return pg
}

func (u *userR) CheckBook(ctx context.Context, book_id string, id float64) pgx.Row {
	pg := u.db.QueryRow(ctx, sql.CheckBookId, book_id, id)
	return pg
}

func (u *userR) TakeBook(ctx context.Context, book_id string, id float64) error {
	_, err := u.db.Exec(ctx, sql.TakeBook, book_id, id)

	return err
}

func (u *userR) GetById(ctx context.Context, id float64) (pgx.Rows, error) {
	pg, err := u.db.Query(ctx, sql.GetById, id)

	return pg, err
}

func (u *userR) ValidateDeleteBook(ctx context.Context, book_id string, id float64) pgx.Row {
	pg := u.db.QueryRow(ctx, sql.CheckUserId, book_id, id)
	
	return pg
}

func (u *userR) DeleteBookById(ctx context.Context, book_id string, id float64) error {
	_, err := u.db.Exec(ctx, sql.DeleteById, book_id, id)

	return err
}

func (u *userR) UserOneBook(ctx context.Context, id float64, book_id string) pgx.Row {
	pg := u.db.QueryRow(ctx, sql.GetOneById, id, book_id)

	return pg
}

func (u *userR) ReqAdmin(ctx context.Context, id float64) error {
	_, err := u.db.Exec(ctx, sql.RequestAdmin, id)

	return err
}

func (u *userR) UpdateRequest(ctx context.Context, answer models.Request, admin_id float64) error {
	_, err := u.db.Exec(ctx, sql.UpdateRequest, answer.UserID, admin_id, answer.Answer)

	return err
}

func (u *userR) TokenReqAdmin(ctx context.Context, id uint64) pgx.Row {
	pg := u.db.QueryRow(ctx, sql.TokenReq, id)

	return pg
}

func (u *userR) CheckReqAdmin(ctx context.Context, id float64) pgx.Row {
	pg := u.db.QueryRow(ctx, sql.CheckReq, id)
	
	return pg
}

func (u *userR) GetAllRequest(ctx context.Context) (pgx.Rows, error) {
	pg, err := u.db.Query(ctx, sql.GetReqAdmin)

	return pg, err
}
func (u *userR) GetCoin(ctx context.Context, id float64) pgx.Row {
	pg := u.db.QueryRow(ctx, sql.GetCoin, id)

	return pg
}


func (u *userR) UpdateCoinUser(ctx context.Context, id float64, coin int) error {
	_, err := u.db.Exec(ctx, sql.UpdateCoin, id, coin - 1)

	return err
}