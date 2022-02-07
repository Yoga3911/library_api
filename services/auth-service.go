package services

import (
	"context"
	"fmt"
	"project_restapi/models"
	"project_restapi/repository"
	"strings"
	"time"
)

type AuthS interface {
	CreateUser(ctx context.Context, user models.Register) error
	VerifyCredential(user models.Login) (string, models.User, error)
	UpdateActive(ctx context.Context, email string) error
}

type authS struct {
	authR repository.AuthR
	jwtS  JWTS
}

func NewAuthS(authR repository.AuthR, jwtS JWTS) AuthS {
	return &authS{authR: authR, jwtS: jwtS}
}

func (a *authS) CreateUser(ctx context.Context, user models.Register) error {
	var (
		name, email int8
	)

	err := a.authR.CheckDuplicate(ctx, user.Name, user.Email).Scan(&name, &email)
	if err != nil {
		return err
	}

	if name != 0 {
		return fmt.Errorf("duplicate name")
	}

	if email != 0 {
		return fmt.Errorf("duplicate email")
	}

	hash, err := hashAndSalt(user.Password)
	if err != nil {
		return err
	}

	err = a.authR.InsertData(ctx, user, hash)
	if err != nil {
		return err
	}

	return nil
}

func (a *authS) VerifyCredential(user models.Login) (string, models.User, error) {
	var usr models.User

	var (
		createT, updateT time.Time
	)
	err := a.authR.VerifyData(user.Email).Scan(&usr.ID, &usr.Name, &usr.Email, &usr.Password, &usr.GenderID, &usr.RoleID, &usr.Coin, &usr.IsActive, &createT, &updateT, &usr.Image)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return "", usr, fmt.Errorf("email not found")
		}
		return "", usr, fmt.Errorf(err.Error())
	}
	usr.CreateAt = createT.Format("02-01-2006")
	usr.UpdateAt = updateT.Format("02-01-2006")

	if !usr.IsActive {
		return "", usr, fmt.Errorf("email not found")
	}

	err = comparePwd([]byte(usr.Password), user.Password)
	if err != nil {
		return "", usr, fmt.Errorf("wrong password")
	}

	token := a.jwtS.GenerateToken(usr.ID, usr.Name, usr.Email, usr.Password, usr.GenderID, usr.RoleID)

	return token, usr, nil
}

func (a *authS) UpdateActive(ctx context.Context, email string) error {
	err := a.authR.UpdateActive(ctx, email)
	
	return err
}