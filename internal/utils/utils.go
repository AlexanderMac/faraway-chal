package utils

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"os"
)

func CheckError(err error) {
	CheckErrorWithMessage(err, "")
}

func CheckErrorWithMessage(err error, message string) {
	if message == "" {
		message = "Fatal error"
	}
	if err != nil {
		fmt.Printf("%s: %s\n", message, err.Error())
		os.Exit(1)
	}
}

func GobRead[T any](r *bufio.Reader) (*T, error) {
	var data T
	decoder := gob.NewDecoder(r)
	err := decoder.Decode(&data)

	return &data, err
}

func GobWrite[T any](w *bufio.Writer, data *T) error {
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(&data)
	if err != nil {
		return err
	}

	err = w.Flush()
	if err != nil {
		return err
	}

	return nil
}
