package app

import (
	"balabanovds/go-stats/app/snmp"
	"balabanovds/go-stats/app/utils"
	"balabanovds/go-stats/app/xmlapi"
	"balabanovds/go-stats/globs"
	"fmt"
	"os"
)

// Fetch elements from servers
func Fetch() []globs.Element {
	nodes := xmlapi.Fetch(globs.Config.Servers)
	if len(nodes) == 0 {
		globs.Fatalf("Fetching from server returned nothing. Please check log file for details.")
	}
	return nodes
}

func All() {
	nodes := Fetch()
	mprs := utils.GetMPRElements(nodes)
	snmp.TrustedManagerForMPRs(mprs)
	utils.Write(nodes)
	fmt.Println("Task done..")
	os.Exit(0)
}

func Encrypt(s string) {
	c := utils.Encrypt(s)
	fmt.Printf("Encrypted: %s\n", c)
	fmt.Println("Task done..")
	os.Exit(0)
}

func SaveData() {
	nodes := Fetch()
	utils.Write(nodes)
	fmt.Println("Task done..")
	os.Exit(0)
}

func CheckTrusted(mprs []globs.Element) {
	snmp.TrustedManagerForMPRs(mprs)
	fmt.Println("Task done..")
	os.Exit(0)
}

func FindDuplicates() {
	nodes := Fetch()
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
