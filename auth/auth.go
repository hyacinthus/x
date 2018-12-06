package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

// ParseJWT 在 echo 的 jwt 插件验证完成后，解析我们需要的数据存储到 context
func ParseJWT(c echo.Context) {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	c.Set("uid", claims["uid"])
	c.Set("oid", claims["oid"])
	c.Set("role", claims["role"])
}
