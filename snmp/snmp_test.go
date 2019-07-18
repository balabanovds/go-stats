package snmp

import (
	"testing"
	g "vimp/globs"
)

func TestTrustedManagerForMPRs(t *testing.T) {
	els := []g.Element{
		{
			IP: "172.21.44.58",
		},
		{
			IP: "10.186.91.9",
		},
	}
	TrustedManagerForMPRs(els)

}
