package main

import (
	"autossl/middleware"
	"autossl/web"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"net/http"
)

func main() {
	e := echo.New()

	// JWT
	//secret := os.Getenv("JWT_SECRET")
	secret := "HBKZ9ec4vrXkbjRIn"
	skipTokenMiddleware := middleware.Skip
	jwtMiddleware := echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(web.JWTCustomClaims)
		},
		SigningKey: []byte(secret),
		Skipper:    skipTokenMiddleware,
	})
	e.Use(jwtMiddleware)

	// Routes
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"status": "UP",
		})
	})
	web.SSLRoutes(e)
	web.LoginRoutes(e)

	// Logger
	e.Use(middleware.RequestLogger())

	// Start the service
	e.Logger.Fatal(e.Start(":1323"))
}
