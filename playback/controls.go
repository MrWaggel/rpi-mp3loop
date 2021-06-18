package playback

import (
	"io/ioutil"
	"mp3loop/easylog"
)

var log = easylog.NewLogModule("Playback")

func PlayFile2(filePath string) error {
	var err error
	var file []byte
	if file, err = ioutil.ReadFile(filePath); err != nil {
		return err
	}
	err = readFile(file)
	if err != nil {
		return err
	}

	return err
}
