package main

import (
	"encoding/json"
	"net/http"

	"github.com/42LoCo42/echotool"
	"github.com/labstack/echo/v4"
)

func GetBooks(c echo.Context) error {
	books := []*Book{}
	if err := Books.FindAll(&books); err != nil {
		return echotool.Die(http.StatusInternalServerError, err, "could not get books")
	}
	return c.JSON(http.StatusOK, books)
}

func GetBook(c echo.Context) error {
	book := &Book{}
	if err := Books.Find(c.Param("id"), book); err != nil {
		return echotool.Die(http.StatusNotFound, err, "book not found")
	}
	return c.JSON(http.StatusOK, book)
}

func PutBook(c echo.Context) error {
	book := &Book{}
	book.ID = c.Param("id")
	Books.Find(book.ID, book)

	if err := json.NewDecoder(c.Request().Body).Decode(book); err != nil {
		return echotool.Die(http.StatusBadRequest, err, "could not decode input")
	}

	if err := Books.Save(book); err != nil {
		return echotool.Die(http.StatusInternalServerError, err, "could not save book")
	}

	return c.JSON(http.StatusOK, book)
}

func DelBook(c echo.Context) error {
	if ok, err := Books.Delete(c.Param("id")); err != nil {
		return echotool.Die(http.StatusInternalServerError, err, "could not delete")
	} else if !ok {
		return echotool.Die(http.StatusNotFound, err, "book not found")
	}
	return c.NoContent(http.StatusOK)
}
