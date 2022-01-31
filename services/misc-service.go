package services

import (
	"fmt"
	// "time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func hashAndSalt(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func comparePwd(hash []byte, pwd string) error {
	err := bcrypt.CompareHashAndPassword(hash, []byte(pwd))

	return err
}

func CheckRole(token string) error {
	t, err := NewJWTS().ValidateToken(token)
	if err != nil {
		return err
	}
	
	claims := t.Claims.(jwt.MapClaims)
	if claims["role_id"] == 1.0 {
		return fmt.Errorf("you are not admin")
	}

	return nil
}

// func dateToMili(ms time.Time) int64 {
//     return ms.UnixMilli()
// }

// func miliToDate(m int64) time.Time {
//     return time.Unix(0, m * int64(time.Millisecond))
// }