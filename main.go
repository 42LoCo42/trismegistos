package main

import (
	"embed"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/42LoCo42/echotool"
	"github.com/albrow/zoom"
)

//go:embed static
var embeddedFS embed.FS

func main() {
	doEmbed := flag.CommandLine.Bool("embed", false, "Embed static data into application")
	flag.Parse()

	var staticFS http.FileSystem = nil
	if *doEmbed {
		staticFS = http.FS(embeddedFS)
	}

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

	e.Start(":8080")
}
