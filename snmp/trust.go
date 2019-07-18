package snmp

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	g "vimp/globs"

	snmp "github.com/soniah/gosnmp"
)

type ipStr string

type managerStatus struct {
	IP     string
	Status string
	Err    error
}

type trustOut struct {
	IP       string
	Statuses []managerStatus
}

func checkTrusted(ip string, chOut chan *trustOut) error {
	g.Debugf(`%s:INFO check host: %v`, tagTrust, ip)
	s := &snmp.GoSNMP{
		Target:    ip,
		Port:      snmpPort,
		Community: snmpCommunity,
		Version:   snmpVersion,
		Timeout:   time.Duration(snmpTimeoutSec) * time.Second,
		Retries:   snmpRetries,
		// Logger:    log.New(os.Stdout, "SNMP", 0),
	}

	err := s.Connect()
	defer s.Conn.Close()
	if err != nil {
		return fmt.Errorf("host: %v err: SNMP Connect() `%v`", s.Target, err)
	}

	res, err := s.WalkAll(opticsIMMgrPollingIPAddress)
	if err != nil {
		return fmt.Errorf("host: %v err: SNMP GetBulk() %v", s.Target, err)
	}

	busyIdx := make(map[int]ipStr)

	for _, p := range res {
		idx, err := getIndex(p.Name)
		if err != nil {
			g.Debugf("%s:ERROR host: %v parsing index of OID %v\n", tagTrust, s.Target, p.Name)
			continue
		}
		ip := ipStr(fmt.Sprintf("%s", p.Value))

		if ip.in(&g.Config.TrustedManagers) || ip.in(&g.Config.Servers) {
			busyIdx[idx] = ip
		}
	}

	freeIdx := getEmptyIndexes(busyIdx)
	mgrs := getManagersToAdd(busyIdx)

	if len(freeIdx) < len(mgrs) {
		return fmt.Errorf("host: %v err: not enough indexes freeIdx %d, managersToAdd %d, busyIdx %+v",
			s.Target, len(freeIdx), len(mgrs), busyIdx)
	}

	out := &trustOut{
		IP: ip,
	}

	for _, ip := range g.Config.TrustedManagers {
		if !ipStr(ip).in(&mgrs) {
			out.Statuses = append(out.Statuses, managerStatus{
				IP:     ip,
				Status: "already trusted",
			})
		}
	}

	for i, ip := range mgrs {
		st := makeTrusted(s, freeIdx[i], string(ip))
		out.Statuses = append(out.Statuses, st)
	}

	chOut <- out
	return nil
}

func getIndex(oid string) (int, error) {
	parts := strings.Split(oid, ".")
	return strconv.Atoi(parts[len(parts)-1])
}

func makeTrusted(s *snmp.GoSNMP, idx int, ip string) managerStatus {
	res := managerStatus{
		IP: ip,
	}

	pdus := []snmp.SnmpPDU{
		{
			Name:  fmt.Sprintf("%s.%d", opticsIMMgrPollingIPAddress, idx),
			Value: ip,
			Type:  snmp.IPAddress,
		},
		{
			Name:  fmt.Sprintf("%s.%d", opticsIMMgrPollingUDPPort, idx),
			Value: []byte(strconv.Itoa(snmpManagerPort)),
			Type:  snmp.OctetString,
		},
		{
			Name:  fmt.Sprintf("%s.%d", opticsIMMgrPollingTimeout, idx),
			Value: snmpPollingTimeout,
			Type:  snmp.Integer,
		},
		{
			Name:  fmt.Sprintf("%s.%d", opticsIMMgrPollingManagerType, idx),
			Value: snmpManagerType,
			Type:  snmp.Integer,
		},
		{
			Name:  fmt.Sprintf("%s.%d", opticsIMMgrPollingRowStatus, idx),
			Value: rowStatusCreateAndGo,
			Type:  snmp.Integer,
		},
	}

	{
		err := deleteTrusted(s, idx)
		if err != nil {
			res.Err = err
			return res
		}
	}
	_, err := s.Set(pdus)
	if err != nil {
		res.Err = fmt.Errorf("host: %v err: SNMP Set() %v", s.Target, err)
		return res
	}
	res.Status = fmt.Sprintf("added at index %d", idx)
	return res
}

func deleteTrusted(s *snmp.GoSNMP, idx int) error {
	_, err := s.Set([]snmp.SnmpPDU{
		{
			Name:  fmt.Sprintf("%s.%d", opticsIMMgrPollingRowStatus, idx),
			Value: rowStatusDestroy,
			Type:  snmp.Integer,
		},
	})

	return err
}

func getEmptyIndexes(busyIdx map[int]ipStr) []int {
	var a []int
	for _, i := range validIndexes {
		if _, ok := busyIdx[i]; !ok {
			a = append(a, i)
		}
	}
	return a
}

func getManagersToAdd(busyIdx map[int]ipStr) []string {
	var a []string
	for _, ip := range g.Config.TrustedManagers {
		found := false
		for _, i := range busyIdx {
			if ip == string(i) {
				found = true
			}
		}
		if !found {
			a = append(a, ip)
		}
	}
	return a
}

func (s ipStr) eq(str string) bool {
	return string(s) == str
}

func (s ipStr) in(arr *[]string) bool {
	in := false
	for _, el := range *arr {
		if string(s) == el {
			in = true
		}
	}
	return in
}
