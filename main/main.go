package main

import (
	"os"
"fmt"
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("example")

// Example format string. Everything except the message has a custom color
// which is dependent on the log level. Many fields have a custom output
// formatting too, eg. the time returns the hour down to the milli second.
var format = logging.MustStringFormatter(
	//`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
	`%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x} %{message}`,
)

// Password is just an example type implementing the Redactor interface. Any
// time this is logged, the Redacted() function will be called.
type Password string

func (p Password) Redacted() interface{} {
	return logging.Redact(string(p))
}

func main() {
	f := os.Stderr
	logFilePath := "./tt.log"
	if _, err := os.Stat(logFilePath); os.IsNotExist(err) {
		f, err = os.Create(logFilePath)
		if err != nil {
			fmt.Println(err.Error)
		}
	}else {
		f, err = os.OpenFile(logFilePath, os.O_APPEND | os.O_WRONLY, 0600)
		if err != nil {
			fmt.Println(err.Error)
		}
	}

	backend := logging.NewLogBackend(f, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	backendLeveled := logging.AddModuleLevel(backendFormatter)
	backendLeveled.SetLevel(logging.ERROR, "")
	logging.SetBackend(backendLeveled)

	log.Debugf("debug %s", Password("secret"))
	log.Info("info")
	log.Notice("notice")
	log.Warning("warning")
	log.Error("err")
	log.Critical("crit")
	defer f.Close()
}
