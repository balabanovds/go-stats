package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	g "vimp/globs"
	"vimp/snmp"
	"vimp/utils"
	"vimp/xmlapi"
)

const (
	version = "0.0.1"
)

func init() {
	exec, err := os.Executable()
	if err != nil {
		g.Fatalf("%v", err)
	}
	g.BaseDir = filepath.Dir(exec)
}

func main() {
	g.InitConst()

	var ips utils.IPFlag
	fAll := flag.Bool("all", false, "1. Fetch all elements from 5620/NSP\n"+
		"2. Run trusted manager for all MPR elements\n"+
		"3. Store elements into files")
	fEncrypt := flag.String("crypt", "", "Encrypt password to store in config sensitive data")
	fStore := flag.Bool("store", false, "Fetch all elements from 5620/NSP and stores into files")
	fTrust := flag.Bool("trust", false, "Fetch all elements and set trusted-managers from config.yml for MPRs.\n"+
		"If use with flag -ip it will run only for mentioned IPs.\n"+
		"If flag -ip not set it will fetch from 5620/NSP")
	fDupl := flag.Bool("dupl", false, "Find duplicate elements on servers")
	fVerb := flag.Bool("verb", false, "Verbose output")
	fVer := flag.Bool("v", false, "Version")
	flag.Var(&ips, "ip", `Comma separated list of IP. Used only with -trust flag`)
	flag.Parse()

	if *fVerb {
		g.DEBUG = true
	}

	if *fAll {
		nodes := fetch()
		mprs := utils.GetMPRElements(nodes)
		snmp.TrustedManagerForMPRs(mprs)
		utils.Write(nodes)
		fmt.Println("Task done..")
		os.Exit(0)
	}

	if len(*fEncrypt) > 0 {
		c := utils.Encrypt(*fEncrypt)
		fmt.Printf("Encrypted: %s\n", c)
		fmt.Println("Task done..")
		os.Exit(0)
	}

	if *fStore {
		nodes := fetch()
		utils.Write(nodes)
		fmt.Println("Task done..")
		os.Exit(0)
	}

	if *fTrust {
		var mprs []g.Element
		if len(ips.GetIPs()) > 0 {
			for _, i := range ips.GetIPs() {
				mprs = append(mprs, g.Element{
					IP: i,
				})
			}
		} else {
			nodes := fetch()
			mprs = utils.GetMPRElements(nodes)
		}
		snmp.TrustedManagerForMPRs(mprs)
		fmt.Println("Task done..")
		os.Exit(0)
	}

	if *fDupl {
		nodes := fetch()
		duplicates := utils.FindDuplicates(nodes)
		fmt.Println("Duplicates found:")
		for ip, els := range duplicates {
			fmt.Printf("%v:\n", ip)
			for _, v := range els {
				fmt.Printf("\t- %+v\n", v)
			}
		}
		os.Exit(0)
	}

	if *fVer {
		fmt.Printf("v%v\n", version)
		os.Exit(0)
	}

	if len(flag.Args()) == 0 {
		flag.Usage()
	}

}

func fetch() []g.Element {
	nodes := xmlapi.Fetch(g.Config.Servers)
	if len(nodes) == 0 {
		g.Fatalf("Fetching from server returned nothing. Please check log file for details.")
	}
	return nodes
}
