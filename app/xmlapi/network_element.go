package xmlapi

import (
	"balabanovds/go-stats/app/utils"
	g "balabanovds/go-stats/globs"
	"encoding/xml"
	"fmt"
	"strings"
)

// Elements contains slice of NetworkElements and possible error
type elements struct {
	Elements []g.Element
	Policy   g.Policy
	Server   string
	Err      error
}

type nodeDiscoveryControlXML struct {
	XMLName xml.Name
	Body    struct {
		XMLName  xml.Name
		Response struct {
			XMLName xml.Name `xml:"findResponse"`
			Result  struct {
				XMLName xml.Name        `xml:"result"`
				Nodes   []nodeDiscovery `xml:"netw.NodeDiscoveryControl"`
			}
		}
	}
}

type nodeDiscovery struct {
	XMLName     xml.Name `xml:"netw.NodeDiscoveryControl"`
	IP          string   `xml:"routerId"`
	ProductType int      `xml:"productType"`
}

func (rx nodeDiscoveryControlXML) entries(server string, policy *g.Policy) []g.Element {
	nodes := rx.Body.Response.Result.Nodes
	m := make([]g.Element, 0, len(nodes))

	for _, item := range nodes {
		e := g.Element{
			IP:          item.IP,
			ElementType: item.ProductType,
			Policy:      *policy,
			Server:      server,
		}
		m = append(m, e)
	}

	return m
}

func fetchElements(server string, policy g.Policy, out chan<- *elements) {
	els := &elements{
		Policy: policy,
		Server: server,
	}
	decrPassword, err := utils.Decrypt(g.Config.Xmlapi.Password)
	if err != nil {
		g.Fatalf("UTILS_CRYPT:FATAL Decrypt password err %v", err)
	}
	query := strings.TrimSpace(fmt.Sprintf(`
	<?xml version="1.0" encoding="UTF-8"?>
	<!--
	get all nodes by MediationPolicyID
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
		<fullClassName>netw.NodeDiscoveryControl</fullClassName>
		<filter>
			<and>
				<equal name="state" value="1"/>
				<equal name="readMediationPolicyId" value="%d"/>
			</and>
		</filter>
		<resultFilter>
					<attribute>routerId</attribute>
					<attribute>productType</attribute>
		</resultFilter>
			</find>
	</SOAP:Body>
	</SOAP:Envelope>
	`, g.Config.Xmlapi.Login, decrPassword, policy.ID))

	resp, err := xmlapiRequest(server, &query)
	if err != nil {
		els.Err = err
		out <- els
		return
	}
	elems, err := unmarshallElements(resp, server, &policy)
	if err != nil {
		els.Err = err
		out <- els
		return
	}
	els.Elements = elems
	out <- els
}

func unmarshallElements(data []byte, server string, policy *g.Policy) ([]g.Element, error) {
	var m nodeDiscoveryControlXML
	err := xml.Unmarshal(data, &m)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall xml responce: %v", err)
	}

	return m.entries(server, policy), nil
}
