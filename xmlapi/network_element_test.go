package xmlapi

import (
	"strings"
	"testing"
	g "vimp/globs"
)

func TestUnmarshallNode(t *testing.T) {
	answer := []byte(strings.TrimSpace(`
	<?xml version="1.0"?>
	<SOAP:Envelope xmlns:SOAP="http://schemas.xmlsoap.org/soap/envelope/">
		<SOAP:Header>
			<header xmlns="xmlapi_1.0">
			<requestID>client1:0</requestID>
			<requestTime>Apr 26, 2019 1:57:49 PM</requestTime>
			<responseTime>Apr 26, 2019 1:57:49 PM</responseTime>
			</header>
		</SOAP:Header>
		<SOAP:Body>
			<findResponse xmlns="xmlapi_1.0">
				<result>
					<netw.NodeDiscoveryControl>
						<routerId>10.72.0.110</routerId>
						<productType>19</productType>
						<children-Set/>
					</netw.NodeDiscoveryControl>
					<netw.NodeDiscoveryControl>
						<routerId>10.72.2.234</routerId>
						<productType>19</productType>
						<children-Set/>
					</netw.NodeDiscoveryControl>
				</result>
			</findResponse>
		</SOAP:Body>
	</SOAP:Envelope>
	`))

	policy := &g.Policy{
		ID: 123,
	}

	nes, _ := unmarshallElements(answer, "1.1.1.1", policy)

	if len(nes) != 2 {
		t.Errorf("Length expected 2, but got %v", len(nes))
	}
}

func TestFetchNEsGorutinesNonBuffered(t *testing.T) {
	g.InitConst()
	ch := make(chan *elements)

	p := g.Policy{
		ID: 96,
	}

	ip := "10.188.33.4"
	expected := 15

	go fetchElements(ip, p, ch)

	nes := <-ch
	if len(nes.Elements) != expected {
		t.Errorf("Get elements length:\n want %d\n got %v", expected, len(nes.Elements))
	}
}
