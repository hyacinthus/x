package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/hyacinthus/x/xerr"
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

// CheckOwner 检查是否机构所有者
func CheckOwner(c echo.Context) error {
	role := int(c.Get("role").(float64))
	if role != 0 {
		return xerr.ErrForbidden
	}
	return nil
}

// CheckAdmin 检查是否机构管理员
func CheckAdmin(c echo.Context) error {
	role := int(c.Get("role").(float64))
	if role > 1 {
		return xerr.ErrForbidden
	}
	return nil
}

// GetUID 获得 user id
func GetUID(c echo.Context) string {
	return c.Get("uid").(string)
}

// GetOID 获得 org id
func GetOID(c echo.Context) string {
	return c.Get("oid").(string)
}
