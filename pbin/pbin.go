package pbin

import (
	"bytes"
	"github.com/Ronmi/pastebin"
	"mp3loop/easylog"
	"mp3loop/settings"
	"os/exec"
	"time"
)

var log = easylog.NewLogModule("pbin")

func Initialize() {
	// check if pbin is set
	k := settings.GetPastebinKey()
	if k == "" {
		return
	}

	uk := settings.GetPastebinUserKey()
	// do a ifconfig
	command := exec.Command("ifconfig")
	output := bytes.Buffer{}
	command.Stdout = &output

	err := command.Run()
	if err != nil {
		log.LogError(err)
		return
	}

	// get body
	body := output.String()

	title := time.Now().Format("2006-01-02 15:04")

	title = "MP3LOOP start " + title
	// make pastebin request
	api := pastebin.API{
		Key: k,
	}

	paste := (&pastebin.Paste{
		Title:      title,
		Content:    body,
		AccessMode: pastebin.Private,
		ExpireAt:   pastebin.In1W,
		Format:     "text",
		UserKey:    uk,
	})

	_, err = api.Post(paste)

	if err != nil {
		log.LogError(err)
		return
	}
}
