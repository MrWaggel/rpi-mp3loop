package playback

import (
	"errors"
	"github.com/hajimehoshi/oto"
	"github.com/tosone/minimp3"
	"mp3loop/fileman"
	"mp3loop/settings"
	"mp3loop/stats"
	"sync"
)

var _stopChannel chan bool
var _controlMutex = sync.Mutex{}

var _running bool
var _runningMutex = sync.Mutex{}

var _volumeMutex = sync.Mutex{}
var _volume = float64(1.0)

func Initialize() {
	// load default song
	defaultfilepath := fileman.DefaultFilePath()
	PlayFile2(defaultfilepath)

	// get init vol
	vol := settings.VolumeSoftwareGet()
	SetVolume(vol)
}

func IsRunning() bool {
	_runningMutex.Lock()
	defer _runningMutex.Unlock()
	return _running
}

func VolumeGetInt() int {
	_volumeMutex.Lock()
	defer _volumeMutex.Unlock()
	return int(_volume * 100)
}

func VolumeGet() float64 {
	_volumeMutex.Lock()
	defer _volumeMutex.Unlock()
	return _volume
}

func SetVolume(vol int) error {
	// Check if not more than 150
	var err error
	if vol > 150 {
		err = errors.New("volume cannot be more than 150")
	} else if vol < 0 {
		err = errors.New("volume cannot be lower than 0")
	}

	if err != nil {
		log.LogError(err)
	}

	_volumeMutex.Lock()
	defer _volumeMutex.Unlock()
	_volume = float64(vol) / 100
	return nil
}

func Stop() {
	if _stopChannel != nil {
		_stopChannel <- true
	}
}

func setRunning(b bool) {
	_runningMutex.Lock()
	defer _runningMutex.Unlock()
	_running = b
}

func ReStart() error {
	var err error
	// Close old one
	Stop()

	_controlMutex.Lock()
	defer _controlMutex.Unlock()
	// Open a new player context
	if _mp3data == nil {
		return nil
	}

	playerContext, err := oto.NewContext(_mp3data.SampleRate, _mp3data.Channels, 2, _bufferSize)
	if err != nil {
		log.LogError(err)
		return err
	}
	player := playerContext.NewPlayer()

	_stopChannel = make(chan bool)

	go func() {
		setRunning(true)
	nestedloop:
		for {
			select {
			case <-_stopChannel:
				_stopChannel = nil
				player.Close()
				playerContext.Close()
				// close
				break nestedloop
			default:
				nextBuf := _mp3reader.BufferNext()
				vol := VolumeGet()
				// Do volume magic
				if false {
					for i := 0; i < len(nextBuf)/2; i++ {
						v16 := int16(nextBuf[2*i]) | (int16(nextBuf[2*i+1]) << 8)
						v16 = int16(float64(v16) * vol)
						nextBuf[2*i] = byte(v16)
						nextBuf[2*i+1] = byte(v16 >> 8)
					}
				}

				if _mp3reader.IsLastChunk() {
					stats.TotalPlaysInc()
					_mp3reader.Reset()
				}

				_, err := player.Write(nextBuf)
				if err != nil {
					log.LogError(err)
				}
			}
		}
		setRunning(false)
	}()
	return nil
}

var _mp3reader *MP3DataReader
var _mp3data *minimp3.Decoder

func readFile(file []byte) error {
	// stop the previous playback (if any)
	Stop()

	// Get basic info
	dec, data, err := minimp3.DecodeFull(file)
	if err != nil {
		log.LogError(err)
		return err
	}

	// Convert byte stream to a custom reader
	_mp3reader = NewMP3DataReader(data)
	_mp3data = dec
	dec.Close()

	return ReStart()
}

type statusData struct {
	Running bool
	Volume  int
}

func Status() statusData {
	return statusData{
		Running: _mp3reader != nil,
		Volume:  int(_volume * 100),
	}
}
