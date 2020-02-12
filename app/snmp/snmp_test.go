package snmp

import (
	g "balabanovds/go-stats/globs"
	"testing"
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
