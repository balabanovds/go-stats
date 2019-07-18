package xmlapi

import (
	"encoding/xml"
	"errors"
	"fmt"
	"strings"
	g "vimp/globs"
	"vimp/utils"
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
	query := strings.TrimSpace(fmt.Sprintf(`
	<?xml version="1.0" encoding="UTF-8"?>
	<!--
	get all Mediation Policy
	-->
	<SOAP:Envelope xmlns:SOAP="http://schemas.xmlsoap.org/soap/envelope/">
	<SOAP:Header>
		<header xmlns="xmlapi_1.0">
		<security>
			<user>%s</user>
			<password hashed="false">%s</password>
		</security>
		<requestID>client1:0</requestID>
		</header>
	</SOAP:Header>
	<SOAP:Body>
		<find xmlns="xmlapi_1.0">
			<fullClassName>security.MediationPolicy</fullClassName>
			<filter>
			</filter>
			<resultFilter>
				<attribute>id</attribute>
				<attribute>securityModel</attribute>
				<attribute>community</attribute>
				<attribute>userName</attribute>
			</resultFilter>
		</find>
	</SOAP:Body>
	</SOAP:Envelope>
	`, g.Config.Xmlapi.Login, decrPassword))

	resp, err := xmlapiRequest(ip, &query)
	if err != nil {
		return &policies{Err: err}
	}

	p := unmarshallMediation(resp)
	if p.Err == nil {
		p.Server = ip
	}
	return p
}

func unmarshallMediation(data []byte) *policies {
	var m mediationXML
	err := xml.Unmarshal(data, &m)
	if err != nil {
		return &policies{
			Err: errors.New("Failed to unmarshall xml responce: " + err.Error()),
		}
	}

	return &policies{
		Policies: m.entries(),
	}
}
