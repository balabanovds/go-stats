package xmlapi

import (
	"balabanovds/go-stats/app/utils"
	"balabanovds/go-stats/globs"
	g "balabanovds/go-stats/globs"
	"bytes"
	"encoding/xml"
	"errors"
)

// policies struct contains slice of fetched policies and error object
type policies struct {
	Policies []g.Policy
	Server   string
	Err      error
}

type mediationXML struct {
	XMLName xml.Name
	Body    struct {
		XMLName  xml.Name
		Response struct {
			XMLName xml.Name `xml:"findResponse"`
			Result  struct {
				XMLName xml.Name    `xml:"result"`
				Policy  []policyXML `xml:"security.MediationPolicy"`
			}
		}
	}
}

type policyXML struct {
	XMLName     xml.Name `xml:"security.MediationPolicy"`
	ID          int      `xml:"id"`
	SnmpVersion string   `xml:"securityModel"`
	Community   string   `xml:"community"`
	UserName    string   `xml:"userName"`
}

func (rx mediationXML) entries() []g.Policy {
	var list []g.Policy

	for _, item := range rx.Body.Response.Result.Policy {
		list = append(list, g.Policy{
			ID:          item.ID,
			SnmpVersion: item.SnmpVersion,
			Community:   item.Community,
			UserName:    item.UserName,
		})
	}

	return list
}

func fetchMediations(ip string, ch chan *policies) {
	ch <- requestMediations(ip)
}

// MediationRun func retrieves from remote server slice of Policy struct
func requestMediations(ip string) *policies {
	decrPassword, err := utils.Decrypt(g.Config.Xmlapi.Password)
	if err != nil {
		g.Fatalf("UTILS_CRYPT:FATAL Decrypt password err %v", err)
	}
	data := struct {
		User     string
		Password string
	}{
		User:     g.Config.Xmlapi.Login,
		Password: decrPassword,
	}

	var query bytes.Buffer
	err = generateText(&query, g.BaseDir, data, "layout", "mediation_policy")
	if err != nil {
		globs.Fatalf("MEDIATION: %v", err)
	}

	resp, err := xmlapiRequest(ip, query.Bytes())
	if err != nil {
		return &policies{Err: err}
	}

	p := unmarshalMediation(resp)
	if p.Err == nil {
		p.Server = ip
	}
	return p
}

func unmarshalMediation(data []byte) *policies {
	var m mediationXML
	err := xml.Unmarshal(data, &m)
	if err != nil {
		return &policies{
			Err: errors.New("Failed to unmarshal xml response: " + err.Error()),
		}
	}

	return &policies{
		Policies: m.entries(),
	}
}
