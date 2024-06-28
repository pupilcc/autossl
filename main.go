package main

import (
	"autossl/infrastructure/acme"
	"autossl/module"
	"autossl/web"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/robfig/cron/v3"
	"log"
	"os"
)

func main() {
	e := echo.New()

	// JWT
	secret := os.Getenv("ADMIN_PASSWORD")
	skipTokenMiddleware := module.Skip
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
	e.Use(module.RequestLogger())

	// Init acme.sh
	acme.InitAcme()
	// Acme Cron
	c := cron.New()
	_, err := c.AddFunc("30 1 * * *", func() {
		err := acme.Cron()
		if err != nil {
			log.Fatal(err)
		}
	})
	if err != nil {
		log.Fatal(err)
	}
	c.Start()

	// Start the service
	e.Logger.Fatal(e.Start(":1323"))
}
