package main

import (
	"encoding/json"
	"net/http"

	"github.com/42LoCo42/echotool"
	"github.com/labstack/echo/v4"
)

func GetBooksRaw() ([]*Book, error) {
	books := []*Book{}
	if err := Books.FindAll(&books); err != nil {
		return nil, echotool.Die(http.StatusInternalServerError, err, "could not get books")
	}
	return books, nil
}

func GetBooks(c echo.Context) error {
	if books, err := GetBooksRaw(); err != nil {
		return err
	} else {
		return c.JSON(http.StatusOK, books)
	}
}

func GetBookRaw(id string) (*Book, error) {
	book := &Book{}
	if err := Books.Find(id, book); err != nil {
		return nil, echotool.Die(http.StatusNotFound, err, "book not found")
	}
	return book, nil
}

func GetBook(c echo.Context) error {
	if book, err := GetBookRaw(c.Param("id")); err != nil {
		return err
	} else {
		return c.JSON(http.StatusOK, book)
	}
}

func PutBookRaw(book *Book) error {
	if err := Books.Save(book); err != nil {
		return echotool.Die(http.StatusInternalServerError, err, "could not save book")
	}
	return nil
}

func PutBook(c echo.Context) error {
	book := &Book{}

	if id := c.Param("id"); id != "" {
		book.ID = id
		Books.Find(book.ID, book)
	}

	if err := json.NewDecoder(c.Request().Body).Decode(book); err != nil {
		return echotool.Die(http.StatusBadRequest, err, "could not decode input")
	}

	if err := PutBookRaw(book); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, book)
}

func DelBookRaw(id string) error {
	if ok, err := Books.Delete(id); err != nil {
		return echotool.Die(http.StatusInternalServerError, err, "could not delete")
	} else if !ok {
		return echotool.Die(http.StatusNotFound, err, "book not found")
	}
	return nil
}

func DelBook(c echo.Context) error {
	if err := DelBookRaw(c.Param("id")); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}
