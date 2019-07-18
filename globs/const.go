package globs

import (
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Configuration struct to parse config.json file
type Configuration struct {
	Servers         []string
	TrustedManagers []string `yaml:"trusted-managers"`

	Xmlapi struct {
		Login    string
		Password string
	}

	DB struct {
		NE       string
		MprStats string `yaml:"mpr-stats"`
		MenStats string `yaml:"men-stats"`
	}

	Snmp struct {
		Threads int
	}

	Log struct {
		Syslog struct {
			Tag      string
			Severity string
			Facility string
		}
	}

	Fs struct {
		DstDir  string `yaml:"destination-dir"`
		NumElem int    `yaml:"number-per-file"`
	}
}

// Config variable that holds all config.yml data
var Config Configuration

// BaseDir is current exec basedir
var BaseDir string

// InitConst init constants
func InitConst() {

	data, err := ioutil.ReadFile(filepath.Join(BaseDir, "config.yml"))
	//data, err := ioutil.ReadFile("/home/dbalaban/dev/go/src/vimp/config.yml")
	if err != nil {
		Fatalf("ReadFile() config err %v", err)
	}
	{
		err := yaml.Unmarshal(data, &Config)
		if err != nil {
			Fatalf("Unmarshal() config err %v", err)
		}
	}
	initLog()
}
