package util

import (
	"gcom-backend/configs"

	"github.com/labstack/echo/v4"
)

func MPMiddleware(mp *configs.MissionPlanner) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("mp", mp)
			return next(c)
		}
	}
}
