package utils

import (
	"bufio"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"io"
	mathRand "math/rand"
	"os"
	"strconv"

	"github.com/AlexanderMac/faraway-chal/internal/constants"
)

func CheckError(err error, prefix string) {
	if prefix == "" {
		prefix = "Fatal error"
	}
	if err != nil {
		fmt.Printf("%s: %s\n", prefix, err.Error())
		os.Exit(1)
	}
}

func IsEOFError(err error) bool {
	return err == io.EOF
}

func ReadMessage[T any](r *bufio.Reader) (*T, error) {
	var ret T
	decoder := gob.NewDecoder(r)
	err := decoder.Decode(&ret)

	return &ret, err
}

func SendMessage[T any](w *bufio.Writer, msgId byte, msg *T) error {
	err := w.WriteByte(msgId)
	if err != nil {
		return err
	}

	encoder := gob.NewEncoder(w)
	err = encoder.Encode(&msg)
	if err != nil {
		return err
	}

	err = w.Flush()
	if err != nil {
		return err
	}

	return nil
}

func CreateChallenge(clientAddr string) string {
	nonce := mathRand.Int()
	data := fmt.Sprintf("%s_%s_%s", constants.SECRET, clientAddr, strconv.Itoa(nonce))

	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}
