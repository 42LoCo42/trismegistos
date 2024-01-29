package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/42LoCo42/echotool"
	"github.com/albrow/zoom"
	"github.com/labstack/echo/v4"
)

func main() {
	pool := zoom.NewPool("localhost:6379")
	defer pool.Close()

	exBook := &Book{
		Title:  "In the Watchful City",
		Author: "S. Qiouyi Lu",
		ISBN10: "1250792983",
		ISBN13: "9781250792983",
	}
	exBook.ModelID()

	exUser := &User{
		Name: "user",
		Hash: echotool.MkHash("user", "pass"),
		Rels: map[string]BookRel{
			exBook.ID: {
				Need: true,
				Done: true,
			},
		},
	}

	books, err := pool.NewCollection(&Book{})
	if err != nil {
		log.Fatal("could not create books collection: ", err)
	}

	if err := books.Save(exBook); err != nil {
		log.Fatal("could not save example book: ", err)
	}

	users, err := pool.NewCollection(&User{})
	if err != nil {
		log.Fatal("could not create users collection: ", err)
	}

	if err := users.Save(exUser); err != nil {
		log.Fatal("could not save example user: ", err)
	}

	e, api := echotool.SetupDefaultEcho(
		staticFS, // directory for static files

		os.Getenv("TRIS_JWT_KEY"), // token key
		"trismegistos",            // token issuer
		time.Hour*24,              // token lifetime

		// find a user by name
		func(name string) (*User, error) {
			user := &User{}
			if err := users.Find(name, user); err != nil {
				return nil, err
			}
			return user, nil
		},

		// get a user's hash
		func(user *User) (string, error) {
			return user.Hash, nil
		},
	)
	if e == nil {
		os.Exit(1)
	}

	api.GET("/self", func(c echo.Context) error {
		return c.JSON(http.StatusOK, c.Get("user"))
	}, echotool.Auth)

	bookG := api.Group("/book")

	bookG.GET("/:id", func(c echo.Context) error {
		book := &Book{}
		if err := books.Find(c.Param("id"), book); err != nil {
			return echotool.Die(http.StatusNotFound, err, "book not found")
		}
		return c.JSON(http.StatusOK, book)
	}, echotool.Auth)

	e.Start(":8080")
}
