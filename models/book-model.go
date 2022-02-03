package models

type Book struct {
	ID       uint64  `json:"id"`
	Title    string  `json:"title"`
	Author   string  `json:"author"`
	Sinopsis string  `json:"sinopsis"`
	Genre    string  `json:"genre"`
	Quantity uint    `json:"quantity"`
	Rating   float64 `json:"rating"`
}

type CUBook struct {
	Title    string `json:"title" validate:"required"`
	Author   string `json:"author" validate:"required"`
	Sinopsis string `json:"sinopsis" validate:"required"`
	GenreID  uint16 `json:"genre_id" validate:"required"`
	Quantity uint   `json:"quantity" validate:"required"`
}

type BookT struct {
	ID       uint64  `json:"id"`
	UserID   uint64  `json:"user_id"`
	BookID   uint64  `json:"book_id"`
	Title    string  `json:"title"`
	Author   string  `json:"author"`
	Sinopsis string  `json:"sinopsis"`
	Genre    string  `json:"genre"`
	Rating   float64 `json:"rating"`
	TDate    string  `json:"transaction_date"`
	RDate    string  `json:"return_date"`
}

type BookReview struct {
	ID       uint64  `json:"id"`
	UserID   uint64  `json:"user_id"`
	BookID   uint64  `json:"book_id"`
	Comment  string  `json:"comment"`
	Rating   float64 `json:"rating"`
	PostDate string  `json:"post_date"`
}

type AddReview struct {
	Comment string `json:"comment"`
	Rating  uint   `json:"rating" validate:"required"`
}
