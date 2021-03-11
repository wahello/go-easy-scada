/*
 * @Descripttion:
 * @version:
 * @Author: Cmpeax
 * @Date: 2019-12-20 10:18:15
 * @LastEditors  : Cmpeax
 * @LastEditTime : 2019-12-20 10:19:31
 */
package midware

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"time"
)

var (
	SIGN_NAME_SCERET = "ESDadminsignKey"
)

// 创建token
func CreateToken(userID string) (string, error) {
	//自定义claim
	claim := jwt.MapClaims{
		"id":  userID,
		"exp": time.Now().Add(time.Hour * time.Duration(72)).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString([]byte(SIGN_NAME_SCERET))
}

// 解析token
func ParseToken(tokenss string) (string, error) {

	token, err := jwt.Parse(tokenss, func(token *jwt.Token) (interface{}, error) {
		return []byte(SIGN_NAME_SCERET), nil
	})
	if err != nil {
		return "", err
	}
	claim, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		err = errors.New("cannot convert claim to mapclaim")
		return "", err
	}

	//验证token，如果token被修改过则为false
	if !token.Valid {
		err = errors.New("token is invalid")
		return "", err
	}

	return claim["id"].(string), nil
}

func CheckToken() gin.HandlerFunc {

	return func(c *gin.Context) {
		// trdb := DB.GetInstance().WebCache.Get()
		// defer trdb.Close()
		isAbort := false
		xtoken := c.Request.Header.Get("X-Token")

		if xtoken != "" {
			_, err := ParseToken(xtoken)
			if err != nil {
				c.AsciiJSON(402, gin.H{
					"message": "没有该token",
					"err":     err.Error(),
				})
				isAbort = true
			}

		} else {
			c.AsciiJSON(402, gin.H{
				"message": "没有携带token",
			})
			isAbort = true
		}
		if isAbort == true {
			c.Abort()
		} else {
			c.Next()
		}

	}
}
