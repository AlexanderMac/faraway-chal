package data

import (
	"bytes"
	"embed"
	"math/rand"
)

//go:embed book.txt
var book embed.FS

func ReadPoem() (string, error) {
	bookData, err := book.ReadFile("book.txt")
	if err != nil {
		return "", err
	}
	lines := bytes.Split(bookData, []byte{'\n'})

	randomLine := rand.Intn(len(lines))
	return string(lines[randomLine]), nil
}
