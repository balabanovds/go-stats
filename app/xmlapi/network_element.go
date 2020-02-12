package xmlapi

import (
	"balabanovds/go-stats/app/utils"
	g "balabanovds/go-stats/globs"
	"bytes"
	"encoding/xml"
	"fmt"
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

	data := struct {
		User     string
		Password string
		PolicyID int
	}{
		User:     g.Config.Xmlapi.Login,
		Password: decrPassword,
		PolicyID: policy.ID,
	}
	var query bytes.Buffer

	err = generateText(&query, g.BaseDir, data, "layout", "network_element")
	if err != nil {
		g.Fatalf(err.Error())
	}

	resp, err := xmlapiRequest(server, query.Bytes())
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
