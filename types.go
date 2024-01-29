package main

import "github.com/albrow/zoom"

type User struct {
	Name string
	Hash string
	Rels map[string]BookRel
}

func (u *User) ModelID() string {
	return u.Name
}

func (u *User) SetModelID(name string) {
	u.Name = name
}

type Book struct {
	zoom.RandomID
	Title  string
	Author string
	ISBN10 string
	ISBN13 string
}

type BookRel struct {
	Need bool
	Owns bool

	Read bool
	Done bool

	BorGive string
	BorFrom string
}
