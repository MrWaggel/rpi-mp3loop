package fileman

import (
	"github.com/tosone/minimp3"
	"io/ioutil"
	"mp3loop/dirs"
	"mp3loop/easylog"
	"mp3loop/settings"
	"os"
	"strings"
	"sync"
)

var lock sync.Mutex
var log = easylog.NewLogModule("FileManager")

func Initialize() {
	ScanFiles()
}

func DefaultFilePath() string {
	filename := settings.GetSelectedFile()
	return dirs.DirFiles() + "/" + filename
}

var cached []*MP3FileInfo

func GetFiles() []*MP3FileInfo {
	return cached
}

func ScanFiles() error {
	files, err := ioutil.ReadDir(dirs.DirFiles())
	if err != nil {
		return err
	}

	out := make([]*MP3FileInfo, 0)
	for _, file := range files {
		if strings.Contains(file.Name(), ".mp3") {
			fin, err := FileInfo(file.Name())
			if err != nil {
				log.LogError(err)
				continue
			}
			out = append(out, &fin)
		}
	}

	lock.Lock()
	defer lock.Unlock()
	cached = out

	return nil
}

func Remove(name string) error {
	abs := dirs.DirFiles() + "/" + name

	err := os.Remove(abs)
	if err != nil {
		log.LogError(err)
	}

	defer ScanFiles()

	return err
}

type MP3FileInfo struct {
	Name       string
	Size       int64
	SampleRate int
	Channels   int
	Kbps       int
	Layer      int
	Selected   bool
}

func (m MP3FileInfo) AbsolutePath() string {
	return dirs.DirFiles() + "/" + m.Name
}

func FileInfo(name string) (MP3FileInfo, error) {
	// get selected file
	selectedFile := settings.GetSelectedFile()

	// get file stats
	abs := dirs.DirFiles() + "/" + name

	i, err := os.Stat(abs)
	if err != nil {
		log.LogError(err)
		return MP3FileInfo{}, err
	}

	out := MP3FileInfo{}
	out.Name = i.Name()
	out.Size = i.Size()

	// decode mp3 for extra juicy details
	b, err := ioutil.ReadFile(abs)
	if err != nil {
		log.LogError(err)
		return MP3FileInfo{}, err
	}

	mi, _, err := minimp3.DecodeFull(b)
	if err != nil {
		log.LogError(err)
		return MP3FileInfo{}, nil
	}

	out.SampleRate = mi.SampleRate
	out.Channels = mi.Channels
	out.Kbps = mi.Kbps
	out.Layer = mi.Layer

	if out.Name == selectedFile {
		out.Selected = true
	} else {
		out.Selected = false
	}

	return out, nil
}

func SaveFile(filename string, b []byte) error {
	// write to disk
	abs := dirs.DirFiles() + "/" + filename
	err := ioutil.WriteFile(abs, b, os.ModePerm)
	if err != nil {
		return err
	}

	// resync files
	return ScanFiles()
}
