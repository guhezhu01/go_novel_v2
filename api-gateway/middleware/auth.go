package middleware

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"time"
)

// MyClaims token生成可以用struct或者map
type RpcClaims struct {
	Id string
	jwt.StandardClaims
}

func authCreateToken(key, id string) (string, bool) {
	// expireTime 设置超时
	expireTime := time.Now().Add(5 * time.Minute).Unix()
	SetClaims := RpcClaims{
		id,
		jwt.StandardClaims{
			ExpiresAt: expireTime,
			Issuer:    viper.GetString("ServiceName"),
		},
	}

	// 根据加密算法和Claims对象来创建Token实例
	reClaim := jwt.NewWithClaims(jwt.SigningMethodHS256, SetClaims)
	// SignedString()利用传入的密钥生成签名字符串
	keyByte := []byte(key)
	token, err := reClaim.SignedString(keyByte)
	if err != nil {
		return "", false
	}
	return token, true
}

func AuthMiddleWare(key, id string) grpc.DialOption {
	return grpc.WithChainUnaryInterceptor(func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		tokenKey, ok := authCreateToken(key, id)
		if !ok {
			log.Println("token create fail,id = ", id)
			return nil
		}
		ctx = metadata.AppendToOutgoingContext(ctx, "token", tokenKey)
		return invoker(ctx, method, req, reply, cc, opts...)
	})
}
