package globs

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
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
func InitConst(file string) {

	data, err := ioutil.ReadFile(file)
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
