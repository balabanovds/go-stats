package snmp

import (
	"log"
	"sync"
	g "vimp/globs"
	"vimp/utils"
)

const (
	tagTrust = "SNMP_TRUST"
)

// TrustedManagerForMPRs starts process to check if ManagerIP
// is included into list of trusted managers
// If not, it will remove the first one in list,
// but not the NSP or SAM IPs
func TrustedManagerForMPRs(elements []g.Element) {
	log.Printf("%s:INFO Start for %d elements\n", tagTrust, len(elements))

	chunks := utils.SplitSlice(elements, g.Config.Snmp.Threads)

	for _, els := range chunks {
		var wg sync.WaitGroup
		wg.Add(len(els))
		chT := make(chan *trustOut)

		for _, e := range els {
			go func(ip string) {
				err := checkTrusted(ip, chT)
				if err != nil {
					g.Debugf("%s:ERROR %v", tagTrust, err)
				}
				wg.Done()
			}(e.IP)
		}

		go func() {
			wg.Wait()
			close(chT)
		}()

		for ch := range chT {
			for _, st := range ch.Statuses {
				g.Debugf("%s:INFO host: %v manager: %v status: %v", tagTrust, ch.IP, st.IP, st.Status)
			}
		}
	}

	log.Printf("%s:INFO Procedure ended. For errors please check log file configured in config.yml.", tagTrust)
}
