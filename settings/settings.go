package settings

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mp3loop/dirs"
	"os"
	"sync"
)

const filename = "settings.json"

var lock sync.Mutex
var sets = new(settings)

func Initialize() {
	fmt.Println(load())
}

func filePath() string {
	return dirs.DirData() + "/settings.json"
}

func load() error {
	b, err := ioutil.ReadFile(filePath())
	if err != nil {
		if err == os.ErrNotExist {
			sets.ForceOutputDeviceOnStart = true
			sets.VolumeSoftware = 50
			sets.VolumeDevice = 100
			return nil
		} else {
			return err
		}
	}

	return json.Unmarshal(b, sets)
}

func save() error {
	lock.Lock()
	defer lock.Unlock()
	m, err := json.Marshal(sets)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filePath(), m, os.ModePerm)
}

type settings struct {
	SelectedFile             string
	DeviceOutputDefaultName  string
	ForceOutputDeviceOnStart bool

	VolumeSoftware int
	VolumeDevice   int
}

func DeviceOutputDefaultGet() string {
	lock.Lock()
	defer lock.Unlock()
	return sets.DeviceOutputDefaultName
}

func DeviceOutputDefaultSet(name string) error {
	lock.Lock()
	sets.DeviceOutputDefaultName = name
	lock.Unlock()
	return save()
}

func VolumeSoftwareGet() int {
	lock.Lock()
	defer lock.Unlock()
	return sets.VolumeSoftware
}
func VolumeSoftwareSet(v int) error {
	lock.Lock()
	sets.VolumeSoftware = v
	lock.Unlock()

	return save()
}

func VolumeDeviceGet() int {
	lock.Lock()
	defer lock.Unlock()
	return sets.VolumeDevice
}

func VolumeDeviceSet(v int) error {
	lock.Lock()
	sets.VolumeDevice = v
	lock.Unlock()

	return save()
}

func ForceOutputDeviceOnStartGet() bool {
	lock.Lock()
	defer lock.Unlock()
	return sets.ForceOutputDeviceOnStart
}

func ForceOutputDeviceOnstartSet(b bool) error {
	lock.Lock()
	sets.ForceOutputDeviceOnStart = b
	lock.Unlock()
	return save()
}

func SetSelectedFile(name string) error {
	sets.SelectedFile = name
	return save()
}

func GetSelectedFile() string {
	lock.Lock()
	defer lock.Unlock()
	return sets.SelectedFile
}
