package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

// ParseJWT 在 echo 的 jwt 插件验证完成后，解析我们需要的数据存储到 context
// TODO: 这个现在只适用于雪豹商情项目，不够通用。
func ParseJWT(c echo.Context) {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	c.Set("uid", claims["uid"])
	c.Set("oid", claims["oid"])
	c.Set("role", claims["role"])
}
