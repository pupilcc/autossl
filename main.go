package main

import (
	"autossl/api"
	"autossl/config"
	"autossl/infrastructure/acme"
	"autossl/infrastructure/database"
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
	skipTokenMiddleware := config.Skip
	jwtMiddleware := echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(api.JWTCustomClaims)
		},
		SigningKey: []byte(secret),
		Skipper:    skipTokenMiddleware,
	})
	e.Use(jwtMiddleware)

	// database
	database.AutoMigrate()

	// Logger
	e.Use(config.RequestLogger())

	// Routes
	api.IndexRoutes(e)
	api.SSLRoutes(e)
	api.LoginRoutes(e)

	// Init acme.sh
	acme.Upgrade()
	acme.ChangeCA()
	acme.UseEmail()
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
