package services

import (
	"context"
	"fmt"
	"project_restapi/models"
	"project_restapi/repository"
	"strings"
)

type AuthS interface {
	CreateUser(ctx context.Context, user models.Register) error
	VerifyCredential(ctx context.Context, user models.Login) (string, models.User, error)
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
	err := a.authR.CheckDuplicate(ctx, user.Name, user.Email)
	if err != nil {
		return err
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

func (a *authS) VerifyCredential(ctx context.Context, user models.Login) (string, models.User, error) {
	usr, err := a.authR.VerifyData(ctx, user.Email)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return "", usr, fmt.Errorf("email not found")
		}
		return "", usr, fmt.Errorf(err.Error())
	}

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