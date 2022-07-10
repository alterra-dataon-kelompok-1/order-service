package middleware

import (
	"strings"

	"github.com/alterra-dataon-kelompok-1/order-service/pkg/response"
	"github.com/alterra-dataon-kelompok-1/order-service/pkg/utils"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type auth struct {
	Roles      map[string]bool
	IDSpecific bool
}

func NewAuthMiddleware(roles map[string]bool, idSpecific bool) *auth {
	return &auth{roles, idSpecific}
}

func (m *auth) Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		authToken := c.Request().Header.Get("Authorization")
		if authToken == "" {
			return response.NewErrorBuilder(response.ErrorConst.Unauthorized).SendJSON(c)
		}

		splitToken := strings.Split(authToken, "Bearer ")

		token, err := utils.VerifyJWT(splitToken[1])
		if err != nil || !token.Valid {
			return response.NewErrorBuilder(response.ErrorConst.Unauthorized).SendJSON(c)
		}

		// var role string
		role := token.Claims.(jwt.MapClaims)["role"].(string)
		val, exist := m.Roles[role]

		if !exist || !val {
			return response.NewErrorBuilder(response.ErrorConst.Unauthorized).SendJSON(c)
		}

		if m.IDSpecific && role != "admin" {
			// Grab uuid from param
			idParam := c.Param("id")

			// Grab uuid from token
			idJwt := token.Claims.(jwt.MapClaims)["id"].(string)

			if idParam != idJwt {
				return response.NewErrorBuilder(response.ErrorConst.Unauthorized).SendJSON(c)
			}
		}
		return next(c)
	}
}

var AdminOnly = map[string]bool{
	"admin":    true,
	"staff":    false,
	"customer": false,
}

var StaffAndAdmin = map[string]bool{
	"admin":    true,
	"staff":    true,
	"customer": false,
}

var SignedIn = map[string]bool{
	"admin":    true,
	"staff":    true,
	"customer": true,
}
