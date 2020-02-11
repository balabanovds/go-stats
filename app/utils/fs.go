package utils

import (
	g "balabanovds/go-stats/globs"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const (
	tagFs = "UTILS_FS"
)

type element struct {
	ip          string
	neType      string
	snmpVersion int
	login       string
}

// Write MPR and MEN elements onto disk
func Write(els []g.Element) {
	if err := os.RemoveAll(filepath.Join(g.BaseDir, g.Config.Fs.DstDir)); err != nil {
		log.Fatal(err)
	}

	m := toServersMap(els)

	for server, elements := range m {
		dirPath := filepath.Join(g.BaseDir, g.Config.Fs.DstDir, server)
		handleDir(dirPath)
		{
			mprs := GetMPRAndMPReElements(elements)
			mprChunks := SplitSlice(mprs, g.Config.Fs.NumElem)
			log.Printf("Found %d MPR elements. Writing to %d files..\n", len(mprs), len(mprChunks))

			counter := 0
			for _, e := range mprChunks {

				filename := filepath.Join(dirPath, fmt.Sprintf("mpr.%d", counter))
				writeToFile(filename, e)
				counter++
			}
		}
		{
			mens := GetMENElements(elements)
			menChunks := SplitSlice(mens, g.Config.Fs.NumElem)
			log.Printf("Found %d MEN elements. Writing to %d files..\n", len(mens), len(menChunks))

			counter := 0
			for _, e := range menChunks {
				filename := filepath.Join(dirPath, fmt.Sprintf("men.%d", counter))
				writeToFile(filename, e)
				counter++
			}
		}

	}

}

func writeToFile(filename string, elements []g.Element) {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		g.Fatalf("%s:FATAL err %v", tagFs, err)
	}
	defer f.Close()
	g.Debugf("%s:INFO Write to file %s, %d elements\n", tagFs, filename, len(elements))

	for _, m := range elements {
		e, _ := getForPerl(&m)
		_, err := fmt.Fprintf(f, "%s,%d,%s\n", e.ip, e.snmpVersion, e.login)
		if err != nil {
			g.Debugf("%s:ERROR Write to file %s err %v", tagFs, filename, err)
		}
	}
}

func handleDir(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			g.Debugf("%s:FATAL MkdirAll() %v", tagFs, err)
			log.Fatal(err)
		}
	}
}
