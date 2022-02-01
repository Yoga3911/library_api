package services

import (
	"context"
	"fmt"
	"log"
	"project_restapi/models"
	"project_restapi/repository"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jackc/pgx/v4/pgxpool"
)

type UserS interface {
	GetAll(ctx context.Context, token string) ([]*models.User, error)
	GetOne(ctx context.Context, token string) (models.User, error)
	Update(ctx context.Context, update models.Update, token string) (string, string, error)
	Delete(ctx context.Context, token string) error
	CheckEmail(ctx context.Context, email string, token string) error
	TakeBook(ctx context.Context, book_id string, token string) error
	GetById(ctx context.Context, token string) ([]*models.BookT, error)
	DeleteBookById(ctx context.Context, book_id string, token string) error
	UserOneBook(ctx context.Context, token string, book_id string) (models.BookT, error)
	ReqAdmin(ctx context.Context, token string) error
	AnsAdmin(ctx context.Context, answer models.Request, token string) (string, error)
	GetAllRequest(ctx context.Context, token string) ([]*models.UserRequest, error)
}

type userS struct {
	userR repository.UserR
	db    *pgxpool.Pool
	jwtS  JWTS
}

func NewUserS(db *pgxpool.Pool, userR repository.UserR, jwtS JWTS) UserS {
	return &userS{db: db, userR: userR, jwtS: jwtS}
}

func (u *userS) GetAll(ctx context.Context, token string) ([]*models.User, error) {
	var users []*models.User

	t, err := u.jwtS.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	claims := t.Claims.(jwt.MapClaims)
	if claims["role_id"] != 2.0 {
		return nil, fmt.Errorf("you are not admin")
	}

	pg, err := u.userR.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	defer pg.Close()

	for pg.Next() {
		var user models.User
		var cDate time.Time
		var uDate time.Time
		err = pg.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.GenderID, &user.RoleID, &user.Coin, &user.IsDeleted, &cDate, &uDate, &user.Image)
		if err != nil {
			log.Println(err)
		}
		user.CreateAt = cDate.Format("02-01-2006")
		user.UpdateAt = uDate.Format("02-01-2006")
		users = append(users, &user)
	}

	return users, nil
}

func (u *userS) GetOne(ctx context.Context, token string) (models.User, error) {
	var user models.User

	t, err := u.jwtS.ValidateToken(token)
	if err != nil {
		return user, err
	}

	claims := t.Claims.(jwt.MapClaims)

	var cDate time.Time
	var uDate time.Time
	err = u.userR.GetOne(ctx, claims["id"]).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.GenderID, &user.RoleID, &user.Coin, &user.IsDeleted, &cDate, &uDate, &user.Image)
	user.CreateAt = cDate.Format("02-01-2006")
	user.UpdateAt = uDate.Format("02-01-2006")
	if err != nil {
		return user, err
	}

	return user, nil
}

func (u *userS) Update(ctx context.Context, update models.Update, t string) (string, string, error) {
	valToken, err := u.jwtS.ValidateToken(t)
	if err != nil {
		return "", "", err
	}

	claims := valToken.Claims.(jwt.MapClaims)

	hash, _ := hashAndSalt(update.Password)

	err = u.userR.Update(ctx, update, claims["id"], hash)
	if err != nil {
		return "", "", err
	}

	token := u.jwtS.GenerateToken(int64(claims["id"].(float64)), update.Name, update.Email, hash, update.GenderID, int16(claims["role_id"].(float64)))

	return token, hash, nil
}

func (u *userS) Delete(ctx context.Context, t string) error {
	valToken, err := u.jwtS.ValidateToken(t)
	if err != nil {
		return err
	}

	claims := valToken.Claims.(jwt.MapClaims)
	err = u.userR.Delete(ctx, claims["id"])

	return err
}

func (u *userS) CheckEmail(ctx context.Context, email string, token string) error {
	var count int
	u.userR.CheckEmail(ctx, email).Scan(&count)
	t, err := u.jwtS.ValidateToken(token)
	if err != nil {
		return err
	}

	claims := t.Claims.(jwt.MapClaims)

	if count > 0 && email != claims["email"] {
		return fmt.Errorf("duplicate email")
	}

	return nil
}

func (u *userS) TakeBook(ctx context.Context, book_id string, token string) error {
	var quantity int
	u.userR.CountQuantity(ctx, book_id).Scan(&quantity)
	if quantity == 0 {
		return fmt.Errorf("out of stock")
	}

	t, err := u.jwtS.ValidateToken(token)
	if err != nil {
		return err
	}

	claims := t.Claims.(jwt.MapClaims)

	var count int
	u.userR.CheckBook(ctx, book_id, claims["id"].(float64)).Scan(&count)
	if count != 0 {
		return fmt.Errorf("you have this book")
	}

	var coin int
	err = u.userR.GetCoin(ctx, claims["id"].(float64)).Scan(&coin)
	if err != nil {
		return err
	}

	if coin == 0 {
		return fmt.Errorf("not enough coin")
	}

	err = u.userR.UpdateCoinUser(ctx, claims["id"].(float64), coin)
	if err != nil {
		return err
	}

	err = u.userR.TakeBook(ctx, book_id, claims["id"])
	if err != nil {
		return err
	}

	return nil
}

func (u *userS) GetById(ctx context.Context, token string) ([]*models.BookT, error) {
	var book []*models.BookT

	t, err := u.jwtS.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	claims := t.Claims.(jwt.MapClaims)

	pg, err := u.userR.GetById(ctx, claims["id"])
	if err != nil {
		return nil, err
	}

	defer pg.Close()
	
	for pg.Next() {
		var b models.BookT
		var tDate time.Time
		var rDate time.Time
		err = pg.Scan(&b.ID, &b.UserID, &b.BookID, &b.Title, &b.Author, &b.Sinopsis, &b.Genre, &b.Rating, &tDate, &rDate)
		if err != nil {
			log.Println(err)
		}
		b.TDate = tDate.Format("02-01-2006")
		b.RDate = rDate.Format("02-01-2006")
		book = append(book, &b)
	}

	return book, nil
}

func (u *userS) DeleteBookById(ctx context.Context, book_id string, token string) error {
	t, err := u.jwtS.ValidateToken(token)
	if err != nil {
		return err
	}

	claims := t.Claims.(jwt.MapClaims)

	var count int
	u.userR.ValidateDeleteBook(ctx, book_id, claims["id"]).Scan(&count)
	if count == 0 {
		return fmt.Errorf("access denied")
	}

	err = u.userR.DeleteBookById(ctx, book_id, claims["id"])
	if err != nil {
		return err
	}

	return nil
}

func (u *userS) UserOneBook(ctx context.Context, token string, book_id string) (models.BookT, error) {
	var book models.BookT

	t, err := u.jwtS.ValidateToken(token)
	if err != nil {
		return book, err
	}

	claims := t.Claims.(jwt.MapClaims)
	var tDate time.Time
	var rDate time.Time
	err = u.userR.UserOneBook(ctx, claims["id"], book_id).Scan(&book.ID, &book.UserID, &book.BookID, &book.Title, &book.Author, &book.Sinopsis, &book.Genre, &book.Rating, &tDate, &rDate)
	book.TDate = tDate.Format("02-01-2006")
	book.RDate = rDate.Format("02-01-2006")
	if err != nil {
		return book, err
	}

	return book, nil
}

func (u *userS) ReqAdmin(ctx context.Context, token string) error {
	t, err := u.jwtS.ValidateToken(token)
	if err != nil {
		return err
	}

	claims := t.Claims.(jwt.MapClaims)

	var count int
	u.userR.CheckReqAdmin(ctx, claims["id"].(float64)).Scan(&count)
	if count != 0 {
		return fmt.Errorf("you already sent a request")
	}

	err = u.userR.ReqAdmin(ctx, claims["id"])
	if err != nil {
		return err
	}

	return nil
}

func (u *userS) AnsAdmin(ctx context.Context, answer models.Request, token string) (string, error) {
	t, err := u.jwtS.ValidateToken(token)
	if err != nil {
		return "", err
	}

	claims := t.Claims.(jwt.MapClaims)

	if claims["role_id"] != 2.0 {
		return "", fmt.Errorf("you are not admin")
	}

	if !answer.Answer {
		err = u.userR.UpdateRequest(ctx, answer, claims["id"])
		if err != nil {
			return "", err
		}

		return "", fmt.Errorf("request declined")
	}

	err = u.userR.UpdateRequest(ctx, answer, claims["id"])
	if err != nil {
		return "", err
	}

	var user models.User

	err = u.userR.TokenReqAdmin(ctx, answer.UserID).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.GenderID, &user.RoleID)
	if err != nil {
		return "", err
	}

	token = u.jwtS.GenerateToken(user.ID, user.Name, user.Email, user.Password, user.GenderID, user.RoleID)

	return token, nil
}

func (u *userS) GetAllRequest(ctx context.Context, token string) ([]*models.UserRequest, error) {
	var users []*models.UserRequest

	t, err := u.jwtS.ValidateToken(token)
	if err != nil {
		return users, err
	}

	claims := t.Claims.(jwt.MapClaims)
	if claims["role_id"] != 2.0 {
		return users, fmt.Errorf("you are not admin")
	}

	pg, err := u.userR.GetAllRequest(ctx)
	if err != nil {
		return users, err
	}

	defer pg.Close()

	for pg.Next() {
		var user models.UserRequest
		var rDate time.Time
		err = pg.Scan(&user.ID, &user.UserID, &rDate)
		user.Request = rDate.Format("02-01-2006")
		if err != nil {
			log.Println(err)
		}
		users = append(users, &user)
	}

	return users, nil
}
