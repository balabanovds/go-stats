package xmlapi

import (
	"bytes"
	"strings"
	"testing"
)

const (
	basedir = "/Users/bds/dev/go/go-stats"
)

func TestGenerateMediationText(t *testing.T) {
	user := struct {
		User     string
		Password string
	}{
		User:     "admin",
		Password: "secret",
	}

	var b bytes.Buffer

	generateText(&b, basedir, user, "layout", "mediation_policy")

	tests := []struct {
		generated string
		expSubstr string
	}{
		{
			generated: b.String(),
			expSubstr: user.User,
		},
		{
			generated: b.String(),
			expSubstr: user.Password,
		},
		{
			generated: b.String(),
			expSubstr: "security.MediationPolicy",
		},
	}

	for _, tt := range tests {
		if !strings.Contains(tt.generated, tt.expSubstr) {
			t.Errorf("expected substring %s not found", tt.expSubstr)
		}
	}
}

func TestGenerateNEText(t *testing.T) {
	user := struct {
		User     string
		Password string
		PolicyID int
	}{
		User:     "admin",
		Password: "secret",
		PolicyID: 111,
	}

	var b bytes.Buffer

	generateText(&b, basedir, user, "layout", "network_element")

	tests := []struct {
		generated string
		expSubstr string
	}{
		{
			generated: b.String(),
			expSubstr: user.User,
		},
		{
			generated: b.String(),
			expSubstr: user.Password,
		},
		{
			generated: b.String(),
			expSubstr: string(user.PolicyID),
		},
		{
			generated: b.String(),
			expSubstr: "netw.NodeDiscoveryControl",
		},
	}

	for _, tt := range tests {
		if !strings.Contains(tt.generated, tt.expSubstr) {
			t.Errorf("expected substring %s not found", tt.expSubstr)
		}
	}
}
