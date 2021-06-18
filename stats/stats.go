package stats

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

var lock sync.Mutex

const _fileName = "stats.json"

func Initialize() {
	// Load from file
	stats = new(statistics)
	load() // ignore eh

	stats.cleanStart()
	stats.lastUpdate = time.Now().Unix()

	// Background updater
	go backgroundUpdater()
}

func save() error {
	lock.Lock()
	defer lock.Unlock()
	// Calculate seconds since
	tn := time.Now().Unix()
	secsSinceLastSave := tn - stats.lastUpdate

	stats.Uptime += secsSinceLastSave
	stats.StartUptime += secsSinceLastSave

	stats.lastUpdate = tn
	// Save to disk

	b, err := json.Marshal(stats)
	if err != nil {
		return err
	}

	// save b to disk
	return ioutil.WriteFile(_fileName, b, os.ModePerm)
}

func load() error {
	lock.Lock()
	defer lock.Unlock()
	// Load file
	b, err := ioutil.ReadFile(_fileName)
	if err != nil {
		if err == os.ErrNotExist {
			return nil
		}
		return err
	}

	// unmarshal
	return json.Unmarshal(b, stats)
}

func backgroundUpdater() {
	for {
		time.Sleep(time.Minute * 2)
		save()
	}
}

type statistics struct {
	// lifetime
	TotalPlays int
	Errors     int
	Warnings   int
	Infos      int
	Uptime     int64

	// Since start
	StartTotalPlays int
	StartErrors     int
	StartWarnings   int
	StartInfos      int
	StartUptime     int64

	// for time calc
	lastUpdate int64
}

func (s *statistics) cleanStart() {
	s.StartTotalPlays = 0
	s.StartErrors = 0
	s.StartWarnings = 0
	s.StartInfos = 0
}

var stats *statistics

func TotalPlaysInc() {
	lock.Lock()
	defer lock.Unlock()
	stats.StartTotalPlays++
	stats.TotalPlays++
}

func TotalErrorsInc() {
	lock.Lock()
	defer lock.Unlock()
	stats.StartErrors++
	stats.Errors++
}

func TotalWarningsInc() {
	lock.Lock()
	defer lock.Unlock()
	stats.Warnings++
	stats.StartWarnings++
}

func TotalInfoInc() {
	lock.Lock()
	defer lock.Unlock()
	stats.Infos++
	stats.StartInfos++
}

func Data() statistics {
	lock.Lock()
	defer lock.Unlock()
	return statistics{
		TotalPlays:      stats.TotalPlays,
		Errors:          stats.Errors,
		Warnings:        stats.Warnings,
		Infos:           stats.Infos,
		Uptime:          stats.Uptime,
		StartTotalPlays: stats.StartTotalPlays,
		StartErrors:     stats.StartErrors,
		StartWarnings:   stats.StartWarnings,
		StartInfos:      stats.StartInfos,
		StartUptime:     stats.StartUptime,
	}
}
