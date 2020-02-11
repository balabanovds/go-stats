package utils

import (
	"fmt"
	"log"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	str := "wert"
	ciph := Encrypt(str)
	fmt.Printf("%v\n", ciph)

	decr, err := Decrypt(ciph)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", decr)
	if str != string(decr) {
		t.Errorf("Wants %s got %s\n", str, decr)
	}
}
