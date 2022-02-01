package sql

const GetAllBook = `SELECT book.id, book.title, book.author, book.sinopsis, genre.genre, book.quantity, book.rating 
					FROM book JOIN genre ON (book.genre_id = genre.id) GROUP BY book.id, genre.genre HAVING (book.quantity > 0)
					ORDER BY id;`

const GetByGenre = `SELECT book.id, book.title, book.author, book.sinopsis, genre.genre, book.quantity, book.rating 
					FROM book JOIN genre ON (book.genre_id = genre.id) WHERE genre.id = $1 
					GROUP BY book.id, genre.genre HAVING (book.quantity > 0) ORDER BY id;`

const AddBook = `INSERT INTO book (title, author, sinopsis, genre_id, quantity, rating) VALUES ($1, $2, $3, $4, $5, 0);`

const UpdateBook = `UPDATE book SET title = $2, author = $3, sinopsis = $4, genre_id = $5, quantity = $6  WHERE id = $1;`

const DeleteBook = `DELETE FROM book WHERE id = $1`

const GetReview = `SELECT * FROM reviews WHERE book_id = $1`

const AddReview = `CALL reviewUser($1, $2, $3, $4)`

const CheckReview = `SELECT COUNT(*) FROM reviews WHERE user_id = $1 AND book_id = $2`

const UpdateReview = `CALL reviewUpdate($1, $2, $3, $4)`

const GetCoin = `SELECT coin FROM users WHERE id = $1`

const UpdateCoin = `UPDATE users SET coin = $2 WHERE id = $1`

const GetGenre = `SELECT genre_id FROM book WHERE id = $1`