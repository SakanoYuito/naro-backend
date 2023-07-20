package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/srinathgs/mysqlstore"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	db   *sqlx.DB
	salt = os.Getenv("HASH_SALT")
)

func main() {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Fatal(err)
	}

	conf := mysql.Config{
		User:      os.Getenv("DB_USERNAME"),
		Passwd:    os.Getenv("DB_PASSWORD"),
		Net:       "tcp",
		Addr:      os.Getenv("DB_HOSTNAME") + ":" + os.Getenv("DB_PORT"),
		DBName:    os.Getenv("DB_DATABASE"),
		ParseTime: true,
		Collation: "utf8mb4_unicode_ci",
		Loc:       jst,
	}

	_db, err := sqlx.Open("mysql", conf.FormatDSN())

	if err != nil {
		log.Fatal(err)
	}

	db = _db

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users (Username VARCHAR(255) PRIMARY KEY, HashedPass VARCHAR(255))")

	if err != nil {
		log.Fatal(err)
	}

	store, err := mysqlstore.NewMySQLStoreFromConnection(db.DB, "sessions", "/", 60*60*24*14, []byte("secret-token"))
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(session.Middleware(store))

	e.POST("/signup", signUpHandler)
	e.POST("/signin", signInHandler)
	e.POST("/signout", signOutHandler)
	e.GET("/ping", func(c echo.Context) error { return c.String(http.StatusOK, "pong") })

	withAuth := e.Group("")
	withAuth.Use(userAuthMiddleware)
	withAuth.GET("/whoami", getWhoAmIHandler)
	withAuth.GET("/cities", getAllCityHandler)
	withAuth.GET("/cities/:cityName", getCityInfoHandler)
	withAuth.GET("/countries", getAllCountryHandler)
	withAuth.GET("/countries/:countryName", getCountryHandler)
	withAuth.GET("/countries/:countryName/cities", getCitiesByCountryHandler)
	withAuth.POST("/cities/newCity", newCityHandler)
	withAuth.PUT("/cities/changeCityInfo/:ID", changeCityInfoHandler)

	e.Start(":8000")
}
