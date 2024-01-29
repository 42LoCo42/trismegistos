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

var (
	Books *zoom.Collection
	Users *zoom.Collection
)

func main() {
	pool := zoom.NewPool("localhost:6379")
	defer pool.Close()

	{
		conn := pool.NewConn()
		conn.Do("flushdb")
		defer conn.Close()
	}

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

	var err error

	Books, err = pool.NewCollectionWithOptions(&Book{}, zoom.DefaultCollectionOptions.WithIndex(true))
	if err != nil {
		log.Fatal("could not create books collection: ", err)
	}

	if err := Books.Save(exBook); err != nil {
		log.Fatal("could not save example book: ", err)
	}

	Users, err = pool.NewCollection(&User{})
	if err != nil {
		log.Fatal("could not create users collection: ", err)
	}

	if err := Users.Save(exUser); err != nil {
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
			if err := Users.Find(name, user); err != nil {
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
	bookG.GET("s", GetBooks, echotool.Auth)
	bookG.GET("/:id", GetBook, echotool.Auth)
	bookG.PUT("/:id", PutBook, echotool.Auth)
	bookG.DELETE("/:id", DelBook, echotool.Auth)

	e.Start(":8080")
}
