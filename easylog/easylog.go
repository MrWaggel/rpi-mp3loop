package easylog

import (
	"encoding/json"
	"io/ioutil"
	"mp3loop/dirs"
	"mp3loop/stats"
	"os"
	"sync"
	"time"
)

const (
	LOG_INFO = iota
	LOG_WARNING
	LOG_ERROR
)

var PrintLogsInConsole = true

const logLimit = 200

var _logCache []LogEntry
var _logLock sync.Mutex

func init() {
	_logCache = make([]LogEntry, 0)
	logsLoad()
}

func GetLog() []string {
	out := make([]string, 0)

	_logLock.Lock()
	defer _logLock.Unlock()

	for _, v := range _logCache {
		s := ""
		// timestamp
		tt := time.Unix(v.Timestamp, 0).Format("2006-01-02 15:04:05")
		s = tt + " " + v.Module

		// type
		t := ""
		switch v.Type {
		case LOG_ERROR:
			t = "ERROR: "
		case LOG_WARNING:
			t = "WARNING: "
		case LOG_INFO:
			t = "INFO: "
		}

		s = s + t + v.Message
		out = append(out, s)
	}
	return out
}

func logfile() string {
	return dirs.DirData() + "/log.json"
}

func logsLoad() error {
	_logLock.Lock()
	defer _logLock.Unlock()
	if _, err := os.Stat(logfile()); err == nil {
		rb, err := ioutil.ReadFile(logfile())
		if err != nil {
			return err
		}
		// path/to/whatever exists
		return json.Unmarshal(rb, &_logCache)
	}
	return nil
}

func logsSave() error {
	jsb, err := json.Marshal(_logCache)
	if err != nil {
		return err
	}

	// write
	return ioutil.WriteFile(logfile(), jsb, os.ModePerm)
}

func logAddEntry(entry LogEntry) {
	_logLock.Lock()
	defer _logLock.Unlock()

	if len(_logCache) >= logLimit {
		// remove last log entry
		_logCache = _logCache[:len(_logCache)-1]
	}

	// add log entry to top
	_logCache = append([]LogEntry{entry}, _logCache...)

	if PrintLogsInConsole {
		str := ""
		switch entry.Type {
		case LOG_INFO:
			str = "INFO: "
		case LOG_WARNING:
			str = "WARNING: "
		case LOG_ERROR:
			str = "ERROR: "
		}

		str = str + entry.Module + " > " + entry.Message
	}
	logsSave()
}

type LogEntry struct {
	Type      int
	Timestamp int64
	Module    string
	Message   string
}

type LogModule struct {
	moduleName string
}

func (lm *LogModule) LogError(err error) {
	if err == nil {
		return
	}

	logentry := LogEntry{}
	logentry.Message = err.Error()
	logentry.Module = lm.moduleName
	logentry.Timestamp = time.Now().Unix()
	logentry.Type = LOG_ERROR

	defer stats.TotalErrorsInc()
	logAddEntry(logentry)
}
func (lm *LogModule) LogWarning(warning string) {

	logentry := LogEntry{}
	logentry.Message = warning
	logentry.Module = lm.moduleName
	logentry.Timestamp = time.Now().Unix()
	logentry.Type = LOG_WARNING
	defer stats.TotalWarningsInc()

	logAddEntry(logentry)
}
func (lm *LogModule) LogInfo(info string) {

	logentry := LogEntry{}
	logentry.Message = info
	logentry.Module = lm.moduleName
	logentry.Timestamp = time.Now().Unix()
	logentry.Type = LOG_INFO

	defer stats.TotalInfoInc()
	logAddEntry(logentry)
}

func NewLogModule(module string) *LogModule {
	o := new(LogModule)
	o.moduleName = module
	return o
}
