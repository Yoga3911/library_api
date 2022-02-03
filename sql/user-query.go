package sql

const GetAll = `SELECT id, name, email, password, gender_id, role_id, coin, is_active, create_at, update_at, image FROM users WHERE is_active = true ORDER BY id`

const GetOne = `SELECT id, name, email, password, gender_id, role_id, coin, is_active, create_at, update_at, image FROM users WHERE id = $1`

const DeleteUser = `UPDATE users SET is_active = false, delete_at = NOW() WHERE id = $1`

const UpdateUser = `UPDATE users SET name = $2, email = $3, gender_id = $4, update_at = NOW(), image = $5 WHERE id = $1`

const TakeBook = `CALL takeBook($1, $2)`

const GetById = `SELECT bt.id, bt.user_id, bt.book_id, b.title, b.author, b.sinopsis, g.genre, b.rating, bt.transaction_date, bt.return_date FROM book_transaction bt JOIN book b ON (bt.book_id=b.id) JOIN genre g ON (g.id=b.genre_id) WHERE user_id = $1`

const GetOneById = `SELECT bt.id, bt.user_id, bt.book_id, b.title, b.author, b.sinopsis, g.genre, b.rating, bt.transaction_date, bt.return_date FROM book_transaction bt JOIN book b ON (bt.book_id=b.id) JOIN genre g ON (g.id=b.genre_id) WHERE user_id = $1 AND book_id = $2`

const DeleteById = `CALL deleteBook($1, $2)`

const CheckBookId = `SELECT COUNT(*) FROM book_transaction WHERE book_id = $1 AND user_id = $2`

const CheckUserId = `SELECT COUNT(*) FROM book_transaction WHERE book_id = $1 AND user_id = $2`

const CheckQuantity = `SELECT quantity FROM book WHERE id = $1`

const RequestAdmin = `INSERT INTO req_admin (user_id, request_date) VALUES ($1, NOW())`

const UpdateRequest = `CALL updateReq($1, $2, $3)`

const TokenReq = `SELECT id, name, email, password, gender_id, role_id FROM users WHERE id = $1`

const CheckReq = `SELECT COUNT(*) FROM req_admin WHERE user_id = $1`

const GetReqAdmin = `SELECT id, user_id, request_date FROM req_admin WHERE is_acc IS NULL`

const CheckPassword = `SELECT COUNT(*) FROM users WHERE id = $1 AND password = $2`

const ChangePass = `UPDATE users SET password = $2 WHERE id = $1`