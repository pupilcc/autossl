package main

import (
	"autossl/middleware"
	"autossl/web"
	"crypto/md5"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	etag "github.com/pablor21/echo-etag/v4"
	"hash"
	"net/http"
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

	// Head request support
	e.Use(AllowHeadRequestsMiddleware())

	// Etag
	e.Use(etag.WithConfig(etag.Config{
		Weak: true,
		HashFn: func(config etag.Config) hash.Hash {
			return md5.New()
		},
	}))

	// Start the service
	e.Logger.Fatal(e.Start(":1323"))
}

func AllowHeadRequestsMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().Method == echo.HEAD {
				c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
				return c.NoContent(http.StatusOK)
			}
			return next(c)
		}
	}
}
