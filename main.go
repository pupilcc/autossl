package main

import (
	"autossl/infrastructure/acme"
	"autossl/middleware"
	"autossl/web"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"os"
)

func main() {
	e := echo.New()

	// JWT
	secret := os.Getenv("ADMIN_PASSWORD")
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
	web.IndexRoutes(e)
	web.SSLRoutes(e)
	web.LoginRoutes(e)

	// Logger
	e.Use(middleware.RequestLogger())

	// Init acme.sh
	acme.InitAcme()

	// Start the service
	e.Logger.Fatal(e.Start(":1323"))
}
