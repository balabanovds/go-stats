package globs

import (
	"log"
	"log/syslog"
)

// DEBUG flag is used for verbose logging to stdout
var DEBUG = false

var (
	l             log.Logger
	severitiesMap = map[string]syslog.Priority{
		"EMERGENCY": syslog.LOG_EMERG,
		"ALERT":     syslog.LOG_ALERT,
		"CRITICAL":  syslog.LOG_CRIT,
		"ERROR":     syslog.LOG_ERR,
		"WARNING":   syslog.LOG_WARNING,
		"NOTICE":    syslog.LOG_NOTICE,
		"INFO":      syslog.LOG_INFO,
		"DEBUG":     syslog.LOG_DEBUG,
	}

	facilitiesMap = map[string]syslog.Priority{
		"KERN":     syslog.LOG_KERN,
		"USER":     syslog.LOG_USER,
		"MAIL":     syslog.LOG_MAIL,
		"DAEMON":   syslog.LOG_DAEMON,
		"AUTH":     syslog.LOG_AUTH,
		"SYSLOG":   syslog.LOG_SYSLOG,
		"LPR":      syslog.LOG_LPR,
		"NEWS":     syslog.LOG_NEWS,
		"UUCP":     syslog.LOG_UUCP,
		"CRON":     syslog.LOG_CRON,
		"AUTHPRIV": syslog.LOG_AUTHPRIV,
		"FTP":      syslog.LOG_FTP,
		"LOCAL0":   syslog.LOG_LOCAL0,
		"LOCAL1":   syslog.LOG_LOCAL1,
		"LOCAL2":   syslog.LOG_LOCAL2,
		"LOCAL3":   syslog.LOG_LOCAL3,
		"LOCAL4":   syslog.LOG_LOCAL4,
		"LOCAL5":   syslog.LOG_LOCAL5,
		"LOCAL6":   syslog.LOG_LOCAL6,
		"LOCAL7":   syslog.LOG_LOCAL7,
	}
)

func initLog() {
	sev := severitiesMap[Config.Log.Syslog.Severity]
	fac := facilitiesMap[Config.Log.Syslog.Facility]

	slg, err := syslog.New(sev|fac, Config.Log.Syslog.Tag)
	if err != nil {
		log.Fatal(err)
	} else {
		l.SetOutput(slg)
	}
}

// Debugf used to print debug info to stdout
func Debugf(format string, a ...interface{}) {
	l.Printf(format, a...)
	if DEBUG {
		log.Printf(format, a...)
	}
}

// Fatalf writes log and exists with err status 1
func Fatalf(format string, a ...interface{}) {
	l.Printf(format, a...)
	log.Fatalf(format, a...)
}
