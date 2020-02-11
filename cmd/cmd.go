package cmd

import (
	"balabanovds/go-stats/app"
	"balabanovds/go-stats/app/utils"
	"balabanovds/go-stats/globs"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

const (
	version = "0.0.1"
)

func Run() {
	exec, err := os.Executable()
	if err != nil {
		globs.Fatalf("%v", err)
	}
	baseDir := filepath.Dir(exec)

	globs.InitConst(filepath.Join(baseDir, "config.yml"))

	if verbose {
		globs.DEBUG = true
	}

	if all {
		app.All()
	}

	if len(crypt) > 0 {
		app.Encrypt(crypt)
	}

	if store {
		app.SaveData()
	}

	if trust {
		var mprs []globs.Element
		allIPs := ips.GetIPs()

		// if we have flag -ip defined than we run only agains theese IPs
		if len(allIPs) > 0 {
			for _, ip := range allIPs {
				mprs = append(mprs, globs.Element{
					IP: ip,
				})
			}
		} else { // we fetch from remote servers
			nodes := app.Fetch()
			mprs = utils.GetMPRElements(nodes)
		}
		app.CheckTrusted(mprs)
	}

	if dupl {
		app.FindDuplicates()
	}

	if v {
		fmt.Printf("v%v\n", version)
		os.Exit(0)
	}

	if len(flag.Args()) == 0 {
		flag.Usage()
	}
}
