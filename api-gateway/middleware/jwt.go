package middleware

import (
	"api-gateway/pkg/errMsg"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"strings"
	"time"
)

var JwtKey = []byte(viper.GetString("JwtKey"))
var code uint32

// MyClaims token生成可以用struct或者map
type MyClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// SetToken 生成token(登录时)
func SetToken(username string) (string, uint32) {
	// expireTime 设置超时
	expireTime := time.Now().Add(10 * time.Hour).Unix()
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
		return "", errMsg.ERROR
	}
	return token, errMsg.SUCCESS
}

// JwtToken 验证token中间件
func JwtToken() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 根据Authorization,获取前端携带的token信息（有两部分，不能直接比较）
		tokenHerder := c.Request.Header.Get("Authorization")
		if tokenHerder == "" {
			code = errMsg.TokenExist
			c.JSON(http.StatusOK, gin.H{
				"status":  code,
				"message": errMsg.GetErrMsg(code),
			})

			// Abort() 结束该函数
			c.Abort()
			return
		}

		// 将tokenHerder分为头部和token两部分
		checkToken := strings.SplitN(tokenHerder, " ", 2)
		if len(checkToken) != 2 && checkToken[0] != "Bearer" {
			code = errMsg.TokenTypeWrong
			c.JSON(http.StatusOK, gin.H{
				"status":  code,
				"message": errMsg.GetErrMsg(code),
			})
			c.Abort()
			return
		}

		// 验证token
		key, typeCode := CheckToken(checkToken[1])
		if typeCode == errMsg.ERROR {
			code = errMsg.TokenWrong
			c.JSON(http.StatusOK, gin.H{
				"status":  code,
				"message": errMsg.GetErrMsg(code),
			})
			c.Abort()
			return
		}
		if time.Now().Unix() > key.ExpiresAt {
			code = errMsg.TokenRuntime
			c.JSON(http.StatusOK, gin.H{
				"status":  code,
				"message": errMsg.GetErrMsg(code),
			})
			c.Abort()
			return
		}

		c.Set("username", key.Username)
		// Next() 继续向下执行
		c.Next()
	}
}

// CheckToken 解析从前端拿到的token
func CheckToken(token string) (*MyClaims, uint32) {
	settoken, _ := jwt.ParseWithClaims(token, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})

	// Claims.(*MyClaims) (断言) 获取到MyClaims结构体类型的内容(包含用户名等信息)
	key, _ := settoken.Claims.(*MyClaims)

	if settoken.Valid {
		return key, errMsg.SUCCESS
	} else {
		return nil, errMsg.ERROR
	}
}
