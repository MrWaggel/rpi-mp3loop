package devices

import (
	"bufio"
	"bytes"
	"errors"
	"mp3loop/easylog"
	"mp3loop/playback"
	"mp3loop/settings"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

var log = easylog.NewLogModule("Devices")

var lock sync.Mutex

func Initialize() {
	// get default device
	log.LogError(getDefaultDevice())
	StartCheck()
}

func StartCheck() {
	// check if we have to force settings
	devname := settings.DeviceOutputDefaultGet()

	// First time run
	if devname == "" {
		devname = _selectedDevice.NameDevice
		defer settings.DeviceOutputDefaultSet(devname)
	}

	if settings.ForceOutputDeviceOnStartGet() {

		// get all devices
		SetDefaultDevice(devname)

		vol := settings.VolumeDeviceGet()

		// force vol
		SetVolume(vol)

	}

}

type DeviceInfo struct {
	IDSink         int
	NameDevice     string
	NameCard       string
	NameAlsaDriver string
	NameAlsaMixer  string
	NameAlsa       string
	NameProduct    string
	Volume         int
	Selected       bool
}

func (di DeviceInfo) GetVolumePercentage() int {
	max := 65536.0
	return int((float64(di.Volume) / max) * 100)
}

var _selectedDevice *DeviceInfo

func GetDevices() ([]*DeviceInfo, error) {
	command := exec.Command("pactl", "list", "sinks")
	output := bytes.Buffer{}
	command.Stdout = &output

	err := command.Run()
	if err != nil {
		log.LogError(err)
		return nil, err
	}

	// Decode data
	scanner := bufio.NewScanner(&output)

	ret := make([]*DeviceInfo, 0)

	var device *DeviceInfo
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "Sink") {
			if device != nil {
				ret = append(ret, device)
			}
			// Split sink
			device = new(DeviceInfo)

			spl := strings.Split(line, "#")

			devID, err := strconv.Atoi(spl[1])
			if err != nil {
				log.LogError(err)
			}

			device.IDSink = devID
			continue
		}

		if strings.Contains(line, "Name: ") {
			spl := strings.Split(line, ": ")
			device.NameDevice = spl[1]
			continue
		}

		if strings.Contains(line, "card_name") {
			spl := strings.Split(line, "= ")
			trim := strings.ReplaceAll(spl[1], "\"", "")
			device.NameCard = trim
			continue
		}

		if strings.Contains(line, "alsa.name") {
			spl := strings.Split(line, "= ")
			trim := strings.ReplaceAll(spl[1], "\"", "")
			device.NameAlsa = trim
			continue
		}

		if strings.Contains(line, "driver.name") {
			spl := strings.Split(line, "= ")
			trim := strings.ReplaceAll(spl[1], "\"", "")
			device.NameAlsaDriver = trim
			continue
		}

		if strings.Contains(line, "product.name") {
			spl := strings.Split(line, "= ")
			trim := strings.ReplaceAll(spl[1], "\"", "")
			device.NameProduct = trim
			continue
		}

		if strings.Contains(line, "alsa.mixer_name") {
			spl := strings.Split(line, "= ")
			trim := strings.ReplaceAll(spl[1], "\"", "")
			device.NameAlsaMixer = trim
			continue
		}

		if strings.Contains(line, "front-left: ") {
			spl := strings.Split(line, "front-left: ")

			spl2 := strings.Split(spl[1], " ")
			valStr := spl2[0]

			v, err := strconv.Atoi(valStr)
			if err != nil {
				log.LogError(err)
			} else {
				device.Volume = v
			}

			continue
		}
	}

	if device != nil {
		ret = append(ret, device)
	}

	if _selectedDevice != nil {
		for _, dev := range ret {
			if dev.IDSink == _selectedDevice.IDSink {
				dev.Selected = true
			} else {
				dev.Selected = false
			}
		}
	}

	lock.Lock()
	defer lock.Unlock()

	return ret, nil
}

func getDefaultDevice() error {
	command := exec.Command("pactl", "info")
	output := bytes.Buffer{}
	command.Stdout = &output

	err := command.Run()
	if err != nil {
		log.LogError(err)
		return err
	}

	// Decode data
	scanner := bufio.NewScanner(&output)
	var sinkname string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "Default Sink") {
			split := strings.Split(line, ": ")
			sinkname = split[1]
		}
	}

	// Get devices
	devs, err := GetDevices()
	if err != nil {
		return err
	}

	// find in devs
	var theDevice *DeviceInfo
	for _, v := range devs {
		if v.NameDevice == sinkname {
			theDevice = v
			break
		}
	}

	if theDevice == nil {
		err := errors.New("the default device returned by 'pactl info' couldn't be find in 'pactl list sinks'")
		log.LogError(err)
		return err
	}

	lock.Lock()
	defer lock.Unlock()
	_selectedDevice = theDevice

	return nil
}

func SetVolume(percentage int) error {
	var err error
	if percentage > 100 {
		err = errors.New("volume percentage cannot be higher than 100")
		log.LogError(err)
		return err
	}

	if percentage < 0 {
		err = errors.New("volume percentage cannot be lower than 0")
		log.LogError(err)
		return err
	}

	// do some magic
	max := 65536
	set := int((float64(max) / 100) * float64(percentage))

	if _selectedDevice == nil {
		err = errors.New("no selected device found, cannot change volume")
		log.LogError(err)
		return err
	}

	// Command
	command := exec.Command("pactl", "set-sink-volume", _selectedDevice.NameDevice, strconv.Itoa(set))
	output := bytes.Buffer{}
	command.Stdout = &output

	err = command.Run()
	if err != nil {
		log.LogError(err)
		return err
	}

	outstr := output.String()

	if len(outstr) > 2 {
		err = errors.New("failed to 'set-sink-volume' for output sink  for " + _selectedDevice.NameDevice + ", error= " + outstr)
		log.LogError(err)
		return err
	}

	return nil
}

func SetDefaultDevice(nameDevice string) error {
	// Stop current playback
	playback.Stop()
	defer playback.ReStart()

	// get current devices
	devs, err := GetDevices()
	if err != nil {
		return err
	}

	var device *DeviceInfo
	for _, v := range devs {
		if v.NameDevice == nameDevice {
			device = v
			break
		}
	}

	if device == nil {
		err = errors.New("cannot change default output: failed to find device with name " + (nameDevice))
		log.LogError(err)
		return err
	}

	// Send command
	// Command
	command := exec.Command("pacmd", "set-default-sink", nameDevice)
	output := bytes.Buffer{}
	command.Stdout = &output

	err = command.Run()
	if err != nil {
		log.LogError(err)
		return err
	}

	outstr := output.String()

	if len(outstr) > 2 {
		err = errors.New("cannot change default output device: " + outstr)
		log.LogError(err)
		return err
	}

	lock.Lock()
	defer lock.Unlock()
	_selectedDevice = device

	return nil
}
