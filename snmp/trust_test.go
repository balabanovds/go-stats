package snmp

import (
	"fmt"
	"testing"
)

func TestCheckTrusted(t *testing.T) {
	ip := "10.186.91.9"
	// ip := "172.21.44.58"
	out := make(chan *trustOut)
	go checkTrusted(ip, out)

	fmt.Printf("%+v", <-out)
}
