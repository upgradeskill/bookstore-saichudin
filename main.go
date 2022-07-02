package main

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	middleware "github.com/labstack/echo/v4/middleware"
	"github.com/saichudin/golang-bookstore/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type App struct {
	*echo.Echo

	books interface {
		All(context.Context) ([]models.Book, error)
		Show(context.Context, uint64) (models.Book, error)
		Create(context.Context, *models.Book) error
		Update(context.Context, uint64, *models.Book) (models.Book, error)
		Delete(context.Context, uint64) error
	}
}

type jwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.StandardClaims
}

func main() {
	db, err := gorm.Open(mysql.Open("root:@tcp(127.0.0.1:8888)/golang_bookstore?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&models.Book{})

	app := &App{
		books: models.BookModel{DB: db},
		Echo:  echo.New(),
	}

	app.SetupRoutes()
	app.Logger.Fatal(app.Start(":3000"))
}

func (app *App) SetupRoutes() {
	app.POST("/login", login)
	app.GET("/books", app.bookIndex)
	app.GET("/books/:id", app.bookShow)

	// Restricted group
	r := app.Group("/admin")

	// Configure middleware with the custom claims type
	config := middleware.JWTConfig{
		Claims:     &jwtCustomClaims{},
		SigningKey: []byte("secret"),
	}
	r.Use(middleware.JWTWithConfig(config))
	r.POST("/books", app.bookCreate)
	r.PUT("/books/:id", app.bookUpdate)
	r.DELETE("/books/:id", app.bookDelete)
}

func (app *App) bookIndex(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*5)
	defer cancel()

	bks, err := app.books.All(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return err
	}

	return c.JSON(http.StatusOK, bks)
}

func (app *App) bookCreate(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*5)
	defer cancel()

	var book models.Book
	if err := c.Bind(&book); err != nil {
		c.JSON(http.StatusBadRequest, "invalid data")
		return err
	}

	if err := app.books.Create(ctx, &book); err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return err
	}

	return c.JSON(http.StatusCreated, book)
}

func (app *App) bookShow(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*5)
	defer cancel()

	paramId := c.Param("id")
	id, _ := strconv.ParseUint(paramId, 10, 64)

	bks, err := app.books.Show(ctx, id)
	if err != nil {
		c.JSON(http.StatusNotFound, map[string]string{"msg": "Not Found"})
		return err
	}

	return c.JSON(http.StatusOK, bks)
}

func (app *App) bookUpdate(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*5)
	defer cancel()

	paramId := c.Param("id")
	id, _ := strconv.ParseUint(paramId, 10, 64)

	var book models.Book
	if err := c.Bind(&book); err != nil {
		c.JSON(http.StatusBadRequest, "invalid data")
		return err
	}

	_, notFound := app.books.Show(ctx, id)
	if notFound != nil {
		c.JSON(http.StatusNotFound, map[string]string{"msg": "Not Found"})
		return notFound
	}

	bks, err := app.books.Update(ctx, id, &book)
	if err != nil {
		c.JSON(http.StatusNotFound, map[string]string{"msg": "Not Found"})
		return err
	}

	return c.JSON(http.StatusOK, bks)
}

func (app *App) bookDelete(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Second*5)
	defer cancel()

	paramId := c.Param("id")
	id, _ := strconv.ParseUint(paramId, 10, 64)

	app.books.Delete(ctx, id)

	return c.JSON(http.StatusOK, map[string]string{"msg": "Data Deleted"})
}

func login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	// Throws unauthorized error
	if username != "admin" || password != "password" {
		return echo.ErrUnauthorized
	}

	// Set custom claims
	claims := &jwtCustomClaims{
		"Admin",
		true,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}
