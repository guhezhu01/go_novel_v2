package middleware

import (
	"comment-service/internal/service"
	"comment-service/pkg/errMsg"
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
)

type RpcClaims struct {
	Id string
	jwt.StandardClaims
}

func checkToken(token, id string) bool {
	setToken, err := jwt.ParseWithClaims(token, &RpcClaims{}, func(token *jwt.Token) (interface{}, error) {
		key := viper.GetString("AuthKey")
		keyByte := []byte(key)
		return keyByte, nil
	})
	if err != nil {
		log.Println(err)
		return false
	}

	key, _ := setToken.Claims.(*RpcClaims)
	if setToken.Valid {
		if key.Id != id {
			return false
		}
	} else {
		return false
	}
	return true
}

func AuthCheckToken(serviceName string) grpc.ServerOption {
	return grpc.ChainUnaryInterceptor(
		func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
			md, _ := metadata.FromIncomingContext(ctx)
			ctxDataToken := md.Get("token")
			ctxDataIds := md.Get(serviceName)

			if len(ctxDataToken) > 0 && len(ctxDataIds) > 0 {
				token := ctxDataToken[0]
				id := ctxDataIds[0]
				ok := checkToken(token, id)
				if !ok {
					return authTokenHandler(ctx, req)
				}
				return handler(ctx, req)
			} else {
				log.Println("no token")
				return authTokenHandler(ctx, req)
			}
			return handler(ctx, req)
		})
}

func authTokenHandler(_ context.Context, _ interface{}) (interface{}, error) {
	resp := &service.CommentsDetailResponse{}
	resp.Code = errMsg.TokenFailed
	resp.Msg = errMsg.GetErrMsg(errMsg.TokenFailed)
	return resp, nil
}
