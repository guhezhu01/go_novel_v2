package middleware

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"time"
)

// MyClaims token生成可以用struct或者map
type MyClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

const JwtKey = "kdfjksfsk"

func CheckToken(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx, status.Error(codes.Unauthenticated, "no authorized")
	}
	tokens := md.Get("token")
	if len(tokens) == 0 {
		return ctx, status.Error(codes.Unauthenticated, "no authorized")
	}

	settoken, err := jwt.ParseWithClaims(tokens[0], &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})
	log.Print(settoken)

	//// Claims.(*MyClaims) (断言) 获取到MyClaims结构体类型的内容(包含用户名等信息)
	//key, _ := settoken.Claims.(*MyClaims)
	//
	//if settoken.Valid {
	//	return key, errMsg.SUCCESS
	//} else {
	//	return nil, errMsg.ERROR
	//}
	return ctx, err
}

// SetToken 生成token(登录时)
func checkToken(username string) (string, bool) {
	// expireTime 设置超时
	expireTime := time.Now().Add(10 * time.Minute).Unix()
	SetClaims := MyClaims{
		username,
		jwt.StandardClaims{
			ExpiresAt: expireTime,
			Issuer:    "gonovel",
		},
	}

	// 根据加密算法和Claims对象来创建Token实例
	reClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, SetClaims)
	// SignedString()利用传入的密钥生成签名字符串
	token, err := reClaim.SignedString(JwtKey)
	if err != nil {
		return "", false
	}
	return token, true
}
