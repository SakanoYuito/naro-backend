package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"strconv"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type City struct {
	ID          int    `json:"id,omitempty"  db:"ID"`
	Name        string `json:"name,omitempty"  db:"Name"`
	CountryCode string `json:"countryCode,omitempty"  db:"CountryCode"`
	District    string `json:"district,omitempty"  db:"District"`
	Population  int    `json:"population,omitempty"  db:"Population"`
}

type Country struct {
	Code    		string 	`json:"code,omitempty"   db:"Code"`
	Name 			string 	`json:"name,omitempty"  db:"Name"`
	Continent		string 	`json:"continent,omitempty"  db:"Continent"`
	Region			string 	`json:"region,omitempty"  db:"Region"`
	SurfaceArea		float64 `json:"surfaceArea,omitempty"  db:"SurfaceArea"`
	IndepYear		int 	`json:"indepYear,omitempty"  db:"IndepYear"`
	Population		int 	`json:"population,omitempty"  db:"Population"`
	LifeExpectancy  float64 `json:"lifeExpectancy,omitempty"  db:"LifeExpectancy"`
	GNP				float64 `json:"gnp,omitempty"  db:"GNP"`
	GNPOld			float64 `json:"gnpOld,omitempty"  db:"GNPOld"`
	LocalName		string 	`json:"localName,omitempty"  db:"LocalName"`
	GovernmentForm	string  `json:"governmentForm,omitempty"  db:"GovernmentForm"`
	HeadOfState		string  `json:"headOfState,omitempty"  db:"HeadOfState"`
	Capital			int 	`json:"capital,omitempty"  db:"Capital"`
	Code2			string  `json:"code2,omitempty"  db:"Code2"`
}

var (
	db *sqlx.DB
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
	fmt.Println("conntected")
	db = _db

	e := echo.New()

	e.GET("/cities/:cityName", getCityInfoHandler)
	e.GET("/countries/:countryName", getCountryHandler)
	e.GET("/cities/allCities", getAllCityHandler)
	e.POST("/cities/newCity", newCityHandler)
	e.PUT("/cities/changeCityInfo/:ID", changeCityInfoHandler)

	e.Start(":8000")
}

func getCityInfoHandler(c echo.Context) error {
	cityName := c.Param("cityName")
	fmt.Println(cityName)

	var city City
	if err := db.Get(&city, "SELECT * FROM city WHERE Name=?", cityName); errors.Is(err, sql.ErrNoRows) {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("No such city Name = %s", cityName))
	} else if err != nil {
		log.Fatalf("DB Error: %s", err)
	}

	return c.JSON(http.StatusOK, city)
}

func getCountryHandler(c echo.Context) error {
	countryName := c.Param("countryName")
	fmt.Println(countryName)

	var country Country
	if err := db.Get(&country, "SELECT * FROM country WHERE Name=?", countryName); errors.Is(err, sql.ErrNoRows) {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("No such city Name = %s", countryName))
	} else if err != nil {
		log.Fatalf("DB Error: %s", err)
	}

	return c.JSON(http.StatusOK, country)
}

func getAllCityHandler(c echo.Context) error {
	var cities []City
	err := db.Select(&cities, "SELECT * FROM city")
	if err != nil {
		log.Fatal(err)
	}
	return c.JSON(http.StatusOK, cities)
}

func newCityHandler(c echo.Context) error {
	city := &City{}
	err := c.Bind(city)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request body")
	}

	_, err = db.Exec("INSERT INTO city (Name, CountryCode, District, Population) VALUES (?, ?, ?, ?);", city.Name, city.CountryCode, city.District, city.Population)

	if err != nil {
		log.Fatalf("DB Error: %s", err) 
	}

	return c.JSON(http.StatusOK, city)

}

func changeCityInfoHandler(c echo.Context) error {
	city := &City{}
	err := c.Bind(city)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request body")
	}
	city.ID, err = strconv.Atoi(c.Param("ID"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "bad request : ID must be integer")
	}
	
	cityBefore := &City{}
	if err := db.Get(cityBefore, "SELECT * FROM city WHERE ID=?", city.ID); errors.Is(err, sql.ErrNoRows) {
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("No such city ID = %d", city.ID))
	} else if err != nil {
		log.Fatalf("DB Error: %s", err)
	}

	if city.Name != "" {
		db.Exec("UPDATE city SET Name = ? WHERE ID = ?;", city.Name, city.ID)
		fmt.Printf("updated Name : %s -> %s\n", cityBefore.Name, city.Name)
		if err != nil {
			log.Fatalf("DB Error: %s", err) 
		}
	}
	if city.CountryCode != "" {
		_, err = db.Exec("UPDATE city SET CountryCode = ? WHERE ID = ?;", city.CountryCode, city.ID)
		fmt.Printf("updated CountryCode : %s -> %s\n", cityBefore.CountryCode, city.CountryCode)
		if err != nil {
			log.Fatalf("DB Error: %s", err) 
		}
	}
	if city.District != "" {
		_, err = db.Exec("UPDATE city SET District = ? WHERE ID = ?;", city.District, city.ID)
		fmt.Printf("updated District : %s -> %s\n", cityBefore.District, city.District)
		if err != nil {
			log.Fatalf("DB Error: %s", err) 
		}
	}
	if city.Population != 0 {
		_, err = db.Exec("UPDATE city SET Population = ? WHERE ID = ?;", city.Population, city.ID)
		fmt.Printf("updated Name : %d -> %d\n", cityBefore.Population, city.Population)
		if err != nil {
			log.Fatalf("DB Error: %s", err) 
		}
	}

	return c.JSON(http.StatusOK, city)

}