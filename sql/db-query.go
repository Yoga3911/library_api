package sql

const Gender = `CREATE TABLE IF NOT EXISTS gender(
	id SERIAL PRIMARY KEY,
	gender VARCHAR(10) NOT NULL
);`

const Role = `CREATE TABLE IF NOT EXISTS user_role(
    id SERIAL PRIMARY KEY,
    role VARCHAR(10) NOT NULL
);`

const Users = `CREATE TABLE IF NOT EXISTS users(
	id BIGSERIAL PRIMARY KEY,
	name VARCHAR(100) NOT NULL,
	email VARCHAR(100) UNIQUE NOT NULL,
	password VARCHAR(100) NOT NULL,
	gender_id SERIAL NOT NULL,
	role_id SERIAL NOT NULL,
	coin INT NOT NULL,
	is_deleted BOOL NOT NULL,
	create_at TIMESTAMP NOT NULL,
	update_at TIMESTAMP NOT NULL,
	delete_at TIMESTAMP,
	image VARCHAR(255) NOT NULL, 
	CONSTRAINT fk_gender FOREIGN KEY (gender_id) REFERENCES gender(id),
	CONSTRAINT fk_role FOREIGN KEY (role_id) REFERENCES user_role(id)
);`

const Genre = `CREATE TABLE IF NOT EXISTS genre(
	id SERIAL PRIMARY KEY,
	genre VARCHAR(20) NOT NULL
);`

const Book = `CREATE TABLE IF NOT EXISTS book(
	id SERIAL PRIMARY KEY,
	title VARCHAR(100) NOT NULL,
	author VARCHAR(100) NOT NULL,
	sinopsis VARCHAR(255) NOT NULL,
	genre_id SERIAL NOT NULL,
	quantity INT NOT NULL,
	rating DECIMAL NOT NULL,
	CONSTRAINT fk_genre FOREIGN KEY (genre_id) REFERENCES genre(id)
);`

const Transaction = `CREATE TABLE IF NOT EXISTS book_transaction(
	id BIGSERIAL PRIMARY KEY,
	user_id SERIAL NOT NULL,
	book_id SERIAL NOT NULL,
	transaction_date TIMESTAMP NOT NULL,
	return_date TIMESTAMP NOT NULL,
	CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id),
	CONSTRAINT fk_book FOREIGN KEY (book_id) REFERENCES book(id) ON DELETE CASCADE
);`

const Req_admin = `CREATE TABLE IF NOT EXISTS req_admin(
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
	admin_id INT,
    is_acc BOOLEAN,
    request_date TIMESTAMP NOT NULL,
    review_date TIMESTAMP,
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_admin FOREIGN KEY (admin_id) REFERENCES users(id)
	);`

const Reviews = `CREATE TABLE IF NOT EXISTS reviews(
	id SERIAL PRIMARY KEY,
	user_id INT NOT NULL,
	book_id INT NOT NULL,
	comment TEXT,
	rating DECIMAL NOT NULL,
	post_date TIMESTAMP NOT NULL,
	CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id),
	CONSTRAINT fk_book FOREIGN KEY (book_id) REFERENCES book(id) ON DELETE CASCADE
);`

const R_gender = `DROP TABLE IF EXISTS gender;`

const R_role = `DROP TABLE IF EXISTS user_role;`

const R_users = `DROP TABLE IF EXISTS users CASCADE;`

const R_genre = `DROP TABLE IF EXISTS genre;`

const R_book = `DROP TABLE IF EXISTS book CASCADE;`

const R_req_admin = `DROP TABLE IF EXISTS req_admin;`

const R_transaction = `DROP TABLE IF EXISTS book_transaction;`

const R_reviews = `DROP TABLE IF EXISTS reviews CASCADE;`

const Func_transaction = `CREATE OR REPLACE PROCEDURE takeBook(bi int, ui int)
LANGUAGE plpgsql
as 
$$
BEGIN
    UPDATE book SET quantity = (SELECT quantity FROM book WHERE id = bi) - 1 WHERE id = bi;

    INSERT INTO book_transaction (user_id, book_id, transaction_date, return_date)
    VALUES (ui, bi, NOW(), NOW() + INTERVAL '30 day');

    COMMIT;
END;
$$;`

const Func_del_transaction = `CREATE OR REPLACE PROCEDURE deleteBook(bi int, ui int)
LANGUAGE plpgsql
as
$$
BEGIN
    UPDATE book SET quantity = (SELECT quantity FROM book WHERE id = bi) + 1 WHERE id = bi;

    DELETE FROM book_transaction WHERE book_id = bi AND user_id = ui;

    COMMIT;
END;
$$;`

const Func_create_validate = `CREATE OR REPLACE FUNCTION createValidate(out count INT, n VARCHAR, e VARCHAR, out nameC INT, out emailC INT)
language plpgsql
as 
$$
BEGIN
	SELECT COUNT(*) INTO nameC FROM users WHERE name =  $2;
	SELECT COUNT(*) INTO emailC FROM users WHERE email = $3;
	SELECT COUNT(*) INTO count FROM users;
END;
$$;`

const Func_update_req = `CREATE OR REPLACE PROCEDURE updateReq(i int, ui int, acc boolean)
LANGUAGE plpgsql
as
$$
    BEGIN
        UPDATE req_admin SET admin_id = $2, is_acc = $3, review_date = NOW() WHERE user_id = $1;
        UPDATE users SET role_id = 2 WHERE id = $1;
    END;
$$;`

const Func_review = `CREATE OR REPLACE PROCEDURE reviewUser(ui int, bi int, com varchar, star int)
LANGUAGE plpgsql
as
$$
DECLARE
    avgR decimal;
BEGIN
    INSERT INTO reviews (user_id, book_id, comment, rating, post_date) VALUES ($1, $2, $3, $4, NOW());
    SELECT AVG(rating) INTO avgR FROM reviews WHERE book_id = $2;
    UPDATE book SET rating = avgR WHERE id = $2;
END;
$$;`

const Func_update_review = `CREATE OR REPLACE PROCEDURE reviewUpdate(ui int, bi int, com varchar, star int)
LANGUAGE plpgsql
as
$$
DECLARE
    avgR decimal;
BEGIN
	UPDATE reviews SET comment = $3, rating = $4, post_date = NOW() WHERE user_id = $1 AND book_id = $2;
    SELECT AVG(rating) INTO avgR FROM reviews WHERE book_id = $2;
    UPDATE book SET rating = avgR WHERE id = $2;
END;
$$;`

const Func_update_coin = `CREATE OR REPLACE PROCEDURE updateCoin()
LANGUAGE plpgsql
AS
$$
DECLARE 
   t_row users%rowtype;
BEGIN
    FOR t_row in SELECT * FROM users LOOP
        if t_row.id NOT IN (SELECT user_id FROM book_transaction WHERE user_id = t_row.id) then
        update users
            set coin = 10
        where id = t_row.id;
        end if;
    END LOOP;
END;
$$;`

const Fetch = `CREATE OR REPLACE PROCEDURE fetchCore()
LANGUAGE plpgsql
as
$$
    BEGIN
        INSERT INTO gender (gender) VALUES ('Laki-laki'), ('Perempuan');
        INSERT INTO genre (genre) VALUES ('Action'), ('Comedy'), ('Romance'), ('Horror');
        INSERT INTO user_role (role) VALUES ('User'), ('Admin');
    END;
$$;`

const Destroy = `CREATE OR REPLACE PROCEDURE destroyCore()
LANGUAGE plpgsql
as
$$
    BEGIN
        DROP PROCEDURE takeBook(int, int);
        DROP PROCEDURE deleteBook(int, int);
        DROP FUNCTION createValidate(out int, varchar, varchar, out int, out int);
		DROP PROCEDURE updateReq(int, int, boolean);
		DROP PROCEDURE reviewUser(int, int, varchar, int);
		DROP PROCEDURE reviewUpdate(int, int, varchar, int);
		DROP PROCEDURE updateCoin();
    END;
$$;`
