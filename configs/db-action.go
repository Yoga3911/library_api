package configs

import (
	"context"
	"project_restapi/sql"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)


func migration(db *pgxpool.Pool) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	_, err := db.Exec(ctx, sql.Gender)
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, sql.Role)
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, sql.Users)
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, sql.Genre)
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, sql.Book)
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, sql.Req_admin)
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, sql.Transaction)
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, sql.Reviews)
	if err != nil {
		return err
	}
	
	_, err = db.Exec(ctx, sql.Func_transaction)
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, sql.Func_del_transaction)
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, sql.Func_create_validate)
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, sql.Func_update_req)
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, sql.Func_review)
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, sql.Func_update_review)
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, sql.Func_update_coin)
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, sql.Fetch)
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, sql.Destroy)
	if err != nil {
		return err
	}

	return nil
}

func rollback(db *pgxpool.Pool) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	_, err := db.Exec(ctx, sql.R_users)
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, sql.R_role)
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, sql.R_gender)
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, sql.R_book)
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, sql.R_genre)
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, sql.R_transaction)
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, sql.R_req_admin)
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, sql.R_reviews)
	if err != nil {
		return err
	}

	return nil
}
