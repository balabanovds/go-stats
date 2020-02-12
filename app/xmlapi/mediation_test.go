package xmlapi

import (
	"strings"
	"testing"
)

func TestMediationUnmarshalling(t *testing.T) {
	answer := []byte(strings.TrimSpace(`
	<?xml version="1.0"?>
	<SOAP:Envelope xmlns:SOAP="http://schemas.xmlsoap.org/soap/envelope/">
		<SOAP:Header>
			<header xmlns="xmlapi_1.0">
			<requestID>client1:0</requestID>
			<requestTime>Apr 29, 2019 6:13:05 PM</requestTime>
			<responseTime>Apr 29, 2019 6:13:05 PM</responseTime>
			</header>
		</SOAP:Header>
		<SOAP:Body>
			<findResponse xmlns="xmlapi_1.0">
				<result>
					<security.MediationPolicy>
						<id>71</id>
						<securityModel>snmpv2c</securityModel>
						<community>nms_snmp</community>
						<userName>N/A</userName>
						<children-Set/>
					</security.MediationPolicy>
					<security.MediationPolicy>
						<id>76</id>
						<securityModel>snmpv2c</securityModel>
						<community>karakuts</community>
						<userName>N/A</userName>
						<children-Set/>
					</security.MediationPolicy>
				</result>
			</findResponse>
		</SOAP:Body>
	</SOAP:Envelope>
	`))

	policies := unmarshalMediation(answer)

	if len(policies.Policies) != 2 {
		t.Errorf("Length expected 2 but got %v", len(policies.Policies))
	}
}

func TestFetchMediationsGoroutine(t *testing.T) {
	ip := "10.188.33.4"

	ch := make(chan *policies)

	go fetchMediations(ip, ch)

	m := <-ch

	if len(m.Policies) != 98 {
		t.Errorf("Fetching mediations expected 98, but got %v", len(m.Policies))
	}

	if m.Server != ip {
		t.Errorf("Fetching mediations expected server %v, but got %v", ip, m.Server)
	}

}

func TestMediation404(t *testing.T) {
	ip := "10.188.33.5"

	policies := requestMediations(ip)
	if policies.Err == nil {
		t.Errorf("Fetching non reachable IP, expected err")
	}

}

func TestMediation404Goroutine(t *testing.T) {
	ip := "10.188.33.5"

	ch := make(chan *policies)

	go fetchMediations(ip, ch)

	m := <-ch

	if m.Err == nil {
		t.Errorf("Expected error")
	}
}
