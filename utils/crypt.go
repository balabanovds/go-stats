package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	g "vimp/globs"
)

const (
	tagCrypt = "UTILS_CRYPT"
	salt     = "I Love Golang!!!"
)

var key string

func init() {
	key = createHash(salt)
}

func createHash(str string) string {
	hasher := md5.New()
	_, err := hasher.Write([]byte(str))
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(hasher.Sum(nil))
}

// Encrypt function is for encrypt strings (passwords)
// to use them in config.json file
func Encrypt(text string) string {
	plaintext := []byte(text)
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		g.Fatalf("%s:FATAL err %v", tagCrypt, err)
	}
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// convert to base64
	return base64.URLEncoding.EncodeToString(ciphertext)
}

// Decrypt function decrypts password to use in further
func Decrypt(str string) (string, error) {
	ciphertext, err := base64.URLEncoding.DecodeString(str)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		g.Fatalf("%s:FATAL err %v", tagCrypt, err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%s", ciphertext), nil
}
