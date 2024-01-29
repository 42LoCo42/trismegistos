package main

import (
	"fmt"
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

	users, err := pool.NewCollection(&User{})
	if err != nil {
		log.Fatal("could not create users collection: ", err)
	}

	if err := users.Save(&User{
		Name: "user",
		Hash: echotool.MkHash("user", "pass"),
	}); err != nil {
		log.Fatal("could not save example user: ", err)
	}

	e, _ := echotool.SetupDefaultEcho(
		staticFS,
		os.Getenv("TRIS_JWT_KEY"),
		"trismegistos",
		time.Hour*24,
		func(name string) (*User, error) {
			user := &User{}
			if err := users.Find(name, user); err != nil {
				return nil, err
			}
			return user, nil
		},
		func(user *User) (string, error) {
			return user.Hash, nil
		},
	)
	if e == nil {
		os.Exit(1)
	}

	e.GET("/foo", func(c echo.Context) error {
		return c.String(http.StatusOK, fmt.Sprint(c.Get("user")))
	})

	e.Start(":8080")
}
