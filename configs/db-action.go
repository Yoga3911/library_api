package configs

import (
	"context"
	"project_restapi/sql"

	"github.com/jackc/pgx/v4/pgxpool"
)


func migration(db *pgxpool.Pool) error {
	_, err := db.Exec(context.Background(), sql.Gender)
	if err != nil {
		return err
	}

	_, err = db.Exec(context.Background(), sql.Role)
	if err != nil {
		return err
	}

	_, err = db.Exec(context.Background(), sql.Users)
	if err != nil {
		return err
	}

	_, err = db.Exec(context.Background(), sql.Genre)
	if err != nil {
		return err
	}

	_, err = db.Exec(context.Background(), sql.Book)
	if err != nil {
		return err
	}

	_, err = db.Exec(context.Background(), sql.Req_admin)
	if err != nil {
		return err
	}

	_, err = db.Exec(context.Background(), sql.Transaction)
	if err != nil {
		return err
	}

	_, err = db.Exec(context.Background(), sql.Reviews)
	if err != nil {
		return err
	}
	
	_, err = db.Exec(context.Background(), sql.Func_transaction)
	if err != nil {
		return err
	}

	_, err = db.Exec(context.Background(), sql.Func_del_transaction)
	if err != nil {
		return err
	}

	_, err = db.Exec(context.Background(), sql.Func_create_validate)
	if err != nil {
		return err
	}

	_, err = db.Exec(context.Background(), sql.Func_update_req)
	if err != nil {
		return err
	}

	_, err = db.Exec(context.Background(), sql.Func_review)
	if err != nil {
		return err
	}

	_, err = db.Exec(context.Background(), sql.Func_update_review)
	if err != nil {
		return err
	}

	_, err = db.Exec(context.Background(), sql.Func_update_coin)
	if err != nil {
		return err
	}

	_, err = db.Exec(context.Background(), sql.Fetch)
	if err != nil {
		return err
	}

	_, err = db.Exec(context.Background(), sql.Destroy)
	if err != nil {
		return err
	}

	return nil
}

func rollback(db *pgxpool.Pool) error {
	_, err := db.Exec(context.Background(), sql.R_users)
	if err != nil {
		return err
	}

	_, err = db.Exec(context.Background(), sql.R_role)
	if err != nil {
		return err
	}

	_, err = db.Exec(context.Background(), sql.R_gender)
	if err != nil {
		return err
	}

	_, err = db.Exec(context.Background(), sql.R_book)
	if err != nil {
		return err
	}

	_, err = db.Exec(context.Background(), sql.R_genre)
	if err != nil {
		return err
	}

	_, err = db.Exec(context.Background(), sql.R_transaction)
	if err != nil {
		return err
	}

	_, err = db.Exec(context.Background(), sql.R_req_admin)
	if err != nil {
		return err
	}

	_, err = db.Exec(context.Background(), sql.R_reviews)
	if err != nil {
		return err
	}

	return nil
}
