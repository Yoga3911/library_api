package services

import (
	"context"
	"fmt"
	"project_restapi/models"
	"project_restapi/repository"
	"strings"
)

type AuthS interface {
	CreateUser(ctx context.Context, user models.Register) (string, string, error)
	VerifyCredential(ctx context.Context, user models.Login) (models.Login, error)
}

type authS struct {
	authR repository.AuthR
	jwtS  JWTS
}

func NewAuthS(authR repository.AuthR, jwtS JWTS) AuthS {
	return &authS{authR: authR, jwtS: jwtS}
}

func (a *authS) CreateUser(ctx context.Context, user models.Register) (string, string, error) {
	var regVal models.RegisterVal

	err := a.authR.QueryCount(ctx, user.Name, user.Email).Scan(&regVal.Count, &regVal.NameC, &regVal.EmailC)
	if err != nil {
		return "", "", err
	}

	if regVal.NameC != 0 {
		return "", "", fmt.Errorf("duplicate name")
	}

	if regVal.EmailC != 0 {
		return "", "", fmt.Errorf("duplicate email")
	}
	
	hash, err := hashAndSalt(user.Password)
	if err != nil {
		return "", "", err
	}

	err = a.authR.InsertData(ctx, user, hash)
	if err != nil {
		return "", "", err
	}

	token := a.jwtS.GenerateToken(regVal.Count+1, user.Name, user.Email, hash, user.GenderID, 1)

	return token, hash, nil
}

func (a *authS) VerifyCredential(ctx context.Context, user models.Login) (models.Login, error) {
	var usr models.Login

	err := a.authR.VerifyData(ctx, user.Email).Scan(&usr.ID, &usr.Name, &usr.Email, &usr.Password, &usr.GenderID, &usr.RoleID, &usr.Coin, &usr.IsDeleted, &usr.Image)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return usr, fmt.Errorf("email not found")
		}
		return usr, fmt.Errorf(err.Error())
	}

	if usr.IsDeleted {
		return usr, fmt.Errorf("user was deleted")
	}

	err = comparePwd([]byte(usr.Password), user.Password)
	if err != nil {
		return usr, fmt.Errorf("wrong password")
	}

	token := a.jwtS.GenerateToken(usr.ID, usr.Name, usr.Email, usr.Password, usr.GenderID, usr.RoleID)
	usr.Token = token

	return usr, nil
}
