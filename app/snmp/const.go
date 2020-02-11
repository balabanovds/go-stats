package snmp

import (
	snmp "github.com/soniah/gosnmp"
)

const (
	//opticsIMMgrPollingInfoTable   = ".1.3.6.1.4.1.637.54.1.1.6.1.8"
	opticsIMMgrPollingIPAddress   = ".1.3.6.1.4.1.637.54.1.1.6.1.8.1.2"
	opticsIMMgrPollingUDPPort     = ".1.3.6.1.4.1.637.54.1.1.6.1.8.1.3"
	opticsIMMgrPollingTimeout     = ".1.3.6.1.4.1.637.54.1.1.6.1.8.1.4"
	opticsIMMgrPollingManagerType = ".1.3.6.1.4.1.637.54.1.1.6.1.8.1.5"
	opticsIMMgrPollingRowStatus   = ".1.3.6.1.4.1.637.54.1.1.6.1.8.1.6"
	//opticsIMMgrPollingIPv6Address = ".1.3.6.1.4.1.637.54.1.1.6.1.8.1.7"

	//rowStatusActive        = 1
	//rowStatusNotInService  = 2
	//rowStatusNotReady      = 3
	rowStatusCreateAndGo = 4
	//rowStatusCreateAndWait = 5
	rowStatusDestroy = 6

	snmpManagerType    = 5
	snmpPort           = 161
	snmpManagerPort    = 162
	snmpPollingTimeout = -1
	snmpVersion        = snmp.Version2c
	snmpCommunity      = "private"
	snmpTimeoutSec     = 3
	snmpRetries        = 5
)

var validIndexes = [...]int{11, 12, 13, 14, 15}
