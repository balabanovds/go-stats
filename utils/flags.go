package utils

import (
	"fmt"
	"strings"
)

// IPFlag struct is used for -ip flag in main func
type IPFlag struct {
	IPs []string
}

// GetIPs return slice of all IPs
func (i *IPFlag) GetIPs() []string {
	return i.IPs
}

func (i *IPFlag) String() string {
	return fmt.Sprint(i.IPs)
}

// Set splits incoming string by comma and sets to IPs
func (i *IPFlag) Set(v string) error {
	if len(i.IPs) > 0 {
		return fmt.Errorf("cannot use ip flag more than once")
	}

	ips := strings.Split(v, ",")
	i.IPs = append(i.IPs, ips...)
	return nil
}
