package utils

import (
	g "balabanovds/go-stats/globs"
	"errors"
)

// GetMPRAndMPReElements filters for ElementType:
// - 18 (MPR-9500)
// - 35 (MPRe)
func GetMPRAndMPReElements(el []g.Element) []g.Element {
	pred := func(e g.Element) bool {
		return e.ElementType == 18 || e.ElementType == 35
	}
	return filter(el, pred)
}

// GetMPRElements filters for ElementType:
// - 18 (MPR-9500)
func GetMPRElements(el []g.Element) []g.Element {
	pred := func(e g.Element) bool {
		return e.ElementType == 18
	}
	return filter(el, pred)
}

// GetMENElements filters for ElementType:
// - 1  (SR-7750)
// - 10 (SAR-7705)
// - 19 (SAS-M-7210)
// - 30 (SAS-D-7210)
// - 54 (SAS-Mxp-7210)
func GetMENElements(el []g.Element) []g.Element {
	pred := func(e g.Element) bool {
		return e.ElementType == 1 ||
			e.ElementType == 10 ||
			e.ElementType == 19 ||
			e.ElementType == 30 ||
			e.ElementType == 54
	}
	return filter(el, pred)
}

// GetOtherElements filters for ElementType:
// - 7  (Generic-Ne)
// - 23 (1830-PSS)
// - 46 (9365-FGW)
func GetOtherElements(el []g.Element) []g.Element {
	pred := func(e g.Element) bool {
		return e.ElementType == 7 ||
			e.ElementType == 23 ||
			e.ElementType == 46
	}
	return filter(el, pred)
}

func filter(elements []g.Element, predicate func(e g.Element) bool) []g.Element {
	var s []g.Element
	for _, e := range elements {
		if predicate(e) {
			s = append(s, e)
		}
	}
	return s
}

func toServersMap(elements []g.Element) map[string][]g.Element {
	m := make(map[string][]g.Element)

	for _, v := range elements {
		m[v.Server] = append(m[v.Server], v)
	}

	return m
}

// FindDuplicates return s a map of elements that present on several servers
func FindDuplicates(elements []g.Element) map[string][]g.Element {
	cntr := make(map[string][]g.Element)
	dupl := make(map[string][]g.Element)
	for _, v := range elements {
		cntr[v.IP] = append(cntr[v.IP], v)
	}

	for k, v := range cntr {
		if len(v) > 1 {
			dupl[k] = v
		}
	}

	return dupl
}

func getForPerl(el *g.Element) (*element, error) {
	e := &element{
		ip: el.IP,
	}
	if el.Policy.SnmpVersion == "snmpv3" {
		e.snmpVersion = 3
		e.login = el.Policy.UserName
	} else {
		e.snmpVersion = 2
		e.login = el.Policy.Community
	}
	e.neType = GetTypeString(el)
	if e.neType == "" {
		return nil, errors.New("unknown NE type")
	}
	return e, nil
}

// GetTypeString just returns MPR or MEN string
func GetTypeString(el *g.Element) string {
	switch el.ElementType {
	case 18, 35:
		return "MPR"
	case 1, 10, 19, 30, 54:
		return "MEN"
	default:
		return "ETC"
	}
}

// SplitSlice does split slice into chunks
func SplitSlice(els []g.Element, chunk int) [][]g.Element {
	var divided [][]g.Element
	for i := 0; i < len(els); i += chunk {
		end := i + chunk
		if end > len(els) {
			end = len(els)
		}
		divided = append(divided, els[i:end])
	}
	return divided
}
