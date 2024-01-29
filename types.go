package main

type User struct {
	Name string
	Hash string
}

func (u *User) ModelID() string {
	return u.Name
}

func (u *User) SetModelID(name string) {
	u.Name = name
}
