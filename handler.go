package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequestBody struct {
	Username string `json:"username,omitempty" form:"username"`
	Password string `json:"password,omitempty" form:"password"`
}

type User struct {
	Username   string `json:"username,omitempty" db:"Username"`
	HashedPass string `json:"hashedpass,omitempty" db:"HashedPass"`
}

type City struct {
	ID          int    `json:"id,omitempty"  db:"ID"`
	Name        string `json:"name,omitempty"  db:"Name"`
	CountryCode string `json:"countryCode,omitempty"  db:"CountryCode"`
	District    string `json:"district,omitempty"  db:"District"`
	Population  int    `json:"population,omitempty"  db:"Population"`
}

type Country struct {
	Code           string  `json:"code,omitempty"   db:"Code"`
	Name           string  `json:"name,omitempty"  db:"Name"`
	Continent      string  `json:"continent,omitempty"  db:"Continent"`
	Region         string  `json:"region,omitempty"  db:"Region"`
	SurfaceArea    float64 `json:"surfaceArea,omitempty"  db:"SurfaceArea"`
	IndepYear      sql.NullInt32     `json:"indepYear,omitempty"  db:"IndepYear"`
	Population     int     `json:"population,omitempty"  db:"Population"`
	LifeExpectancy sql.NullFloat64 `json:"lifeExpectancy,omitempty"  db:"LifeExpectancy"`
	GNP            sql.NullFloat64 `json:"gnp,omitempty"  db:"GNP"`
	GNPOld         sql.NullFloat64 `json:"gnpOld,omitempty"  db:"GNPOld"`
	LocalName      string  `json:"localName,omitempty"  db:"LocalName"`
	GovernmentForm string  `json:"governmentForm,omitempty"  db:"GovernmentForm"`
	HeadOfState    sql.NullString  `json:"headOfState,omitempty"  db:"HeadOfState"`
	Capital        sql.NullInt32     `json:"capital,omitempty"  db:"Capital"`
	Code2          string  `json:"code2,omitempty"  db:"Code2"`
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

func getAllCountryHandler(c echo.Context) error {
	var countries []Country
	err := db.Select(&countries, "SELECT * FROM country")
	if err != nil {
		log.Fatal(err)
	}
	return c.JSON(http.StatusOK, countries)
}

func getCitiesByCountryHandler(c echo.Context) error {
	var cities []City
	countryName := c.Param("countryName")
	var CC string
	if err := db.Get(&CC, "SELECT Code FROM country WHERE Name=?", countryName); errors.Is(err, sql.ErrNoRows){
		return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("No such country Name = %s", countryName))
	} else if err != nil {
		log.Fatalf("DB Error: %s", err)
	}

	err := db.Select(&cities, "SELECT * FROM city WHERE CountryCode=?", CC)
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

func signUpHandler(c echo.Context) error {
	var req LoginRequestBody
	c.Bind(&req)

	if req.Password == "" || req.Username == "" {
		return c.String(http.StatusBadRequest, "Username or Password if empty")
	}

	var count int
	err := db.Get(&count, "SELECT COUNT(*) FROM users WHERE Username=?", req.Username)
	if err != nil {
		log.Println(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	if count > 0 {
		return c.String(http.StatusConflict, "Username is already used")
	}

	pw := req.Password + salt
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	_, err = db.Exec("INSERT INTO users (Username, HashedPass) VALUES (?, ?)", req.Username, hashedPass)
	if err != nil {
		log.Println(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusCreated)
}

func signInHandler(c echo.Context) error {
	var req LoginRequestBody
	c.Bind(&req)

	if req.Password == "" || req.Username == "" {
		return c.String(http.StatusBadRequest, "Username or Password is empty")
	}
	user := User{}
	err := db.Get(&user, "SELECT * FROM users WHERE username=?", req.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return c.NoContent(http.StatusUnauthorized)
		} else {
			log.Println(err)
			return c.NoContent(http.StatusInternalServerError)
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPass), []byte(req.Password+salt))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return c.NoContent(http.StatusUnauthorized)
		} else {
			return c.NoContent(http.StatusInternalServerError)
		}
	}

	sess, err := session.Get("sessions", c)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusInternalServerError, "something wrong in getting session")
	}
	sess.Values["userName"] = req.Username
	sess.Save(c.Request(), c.Response())

	return c.NoContent(http.StatusOK)
}

func signOutHandler(c echo.Context) error {
	sess, err := session.Get("sessions", c)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusInternalServerError, "something wrong in getting session")
	}
	sess.Options.MaxAge = -1
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		log.Fatal("Failed to delete session", err)
	}
	return c.NoContent(http.StatusOK)
}

func userAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get("sessions", c)
		if err != nil {
			fmt.Println(err)
			return c.String(http.StatusInternalServerError, "something wrong in getting session")
		}
		if sess.Values["userName"] == nil {
			return c.String(http.StatusUnauthorized, "please singin")
		}
		c.Set("userName", sess.Values["userName"].(string))

		return next(c)
	}
}


func getWhoAmIHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, struct{ Username string }{
		Username: c.Get("userName").(string),
	})
}
