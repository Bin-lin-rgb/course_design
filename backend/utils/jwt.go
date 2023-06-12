package utils

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type JWT struct {
}

var jwtSecret []byte

func init() {
	jwtSecret = []byte("tce_label")
}

// Claims （记录用户实体）
type Claims struct {
	UserID   int64  `json:"userID"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// TokenNew 生成token
func (j *JWT) TokenNew(userID int64, username string) (string, error) {
	now := time.Now()
	expireTime := now.Add(12 * time.Hour)
	claims := Claims{
		UserID:   userID,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "localhost",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//用密钥进行签名
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}

// TokenParse 解析token
func (j *JWT) TokenParse(token string) (*Claims, error) {
	// 进行解析鉴权声明
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
