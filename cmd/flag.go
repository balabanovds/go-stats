package cmd

import (
	"balabanovds/go-stats/app/utils"
	"flag"
)

var (
	all     bool
	crypt   string
	store   bool
	trust   bool
	dupl    bool
	verbose bool
	v       bool
	ips     utils.IPFlag
)

func init() {
	flag.BoolVar(&all, "all", false, "1. Fetch all elements from 5620/NSP\n"+
		"2. Run trusted manager for all MPR elements\n"+
		"3. Store elements into files")
	flag.StringVar(&crypt, "crypt", "", "Encrypt password to store in config sensitive data")
	flag.BoolVar(&store, "store", false, "Fetch all elements from 5620/NSP and stores into files")
	flag.BoolVar(&trust, "trust", false, "Fetch all elements and set trusted-managers from config.yml for MPRs.\n"+
		"If use with flag -ip it will run only for mentioned IPs.\n"+
		"If flag -ip not set it will fetch from 5620/NSP")
	flag.BoolVar(&dupl, "dupl", false, "Find duplicate elements on servers")
	flag.BoolVar(&verbose, "verb", false, "Verbose output")
	flag.BoolVar(&v, "v", false, "Version")

	flag.Var(&ips, "ip", `Comma separated list of IP. Used only with -trust flag`)
	flag.Parse()
}
