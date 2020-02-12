package xmlapi

import (
	g "balabanovds/go-stats/globs"
	"log"
	// g "vimp/globs"
)

const (
	tagHTTP      = "XMLAPI_HTTP"
	tagMediation = "XMLAPI_MEDIATION"
	tagElements  = "XMLAPI_ELEMENTS"
)


// Fetch function fetches mediations and network_elements
// from remote servers via XMLAPI interface
func Fetch(ipList []string) []g.Element {
	log.Printf("Fetching elements from remote servers %+v\n", g.Config.Servers)

	chP := make(chan *policies)
	chE := make(chan *elements)

	for _, ip := range ipList {
		go fetchMediations(ip, chP)
	}

	counter := 0

	for i := 0; i < len(ipList); i++ {
		ch := <-chP
		if ch.Err != nil {
			g.Debugf("%s:ERROR Fetching %v", tagMediation, ch.Err)
			continue
		}
		policies := ch.Policies
		counter += len(policies)
		for _, p := range policies {
			go fetchElements(ch.Server, p, chE)
		}
	}

	els := make([]g.Element, 0)

	//for e := range chE {
	//	if len(e.Elements) == 0 {
	//		g.Debugf("%s:INFO Elements are 0 at server %v policyId %v. Skipping..\n",
	//			tagElements, e.Server, e.Policy.ID)
	//		continue
	//	}
	//	if e.Err != nil {
	//		g.Debugf("%s:ERROR Fetching server %v policy %v: %v\n",
	//			tagElements, e.Server, e.Policy.ID, e.Err.Error())
	//		continue
	//	}
	//	g.Debugf("%s:INFO found %d elements for policyId %d server %s",
	//		tagElements, len(e.Elements), e.Policy.ID, e.Server)
	//
	//	els = append(els, e.Elements...)
	//}

	for i := 0; i < counter; i++ {
		e := <-chE
		if len(e.Elements) == 0 {
			g.Debugf("%s:INFO Elements are 0 at server %v policyId %v. Skipping..\n",
				tagElements, e.Server, e.Policy.ID)
			continue
		}
		if e.Err != nil {
			g.Debugf("%s:ERROR Fetching server %v policy %v: %v\n",
				tagElements, e.Server, e.Policy.ID, e.Err.Error())
			continue
		}
		g.Debugf("%s:INFO found %d elements for policyId %d server %s",
			tagElements, len(e.Elements), e.Policy.ID, e.Server)

		els = append(els, e.Elements...)
	}

	return els
}
