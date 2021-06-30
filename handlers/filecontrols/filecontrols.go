package filecontrols

import (
	"bytes"
	"errors"
	"github.com/kr/pretty"
	"github.com/labstack/echo"
	"github.com/tosone/minimp3"
	"io"
	"mp3loop/easylog"
	"mp3loop/fileman"
	"mp3loop/handlers"
	"mp3loop/playback"
	"mp3loop/settings"
)

var log = easylog.NewLogModule("WebFileHandler")

const maxFileSize = 1000 * 1000 * 15

func Initialize(e *echo.Echo) {
	e.GET("/files/all", GetAll)
	e.POST("/files/info", FileInfo)
	e.POST("/files/delete", RemoveFile)
	e.POST("/files/use", FileSetDefault)
	e.POST("/files/add", AddFile)
}

func GetAll(c echo.Context) error {
	files := fileman.GetFiles()
	return c.JSON(200, files)
}
func AddFile(c echo.Context) error {
	// get the file

	ffile, err := c.FormFile("file")
	if err != nil {
		log.LogError(err)
		return c.JSON(200, handlers.ResponseError(err))
	}

	if ffile.Size > maxFileSize {
		err = errors.New("mp3 file upload can only be 15 megabytes")
		log.LogWarning("an upload bigger than 15 megabytes was attempted")
		return c.JSON(200, handlers.ResponseError(err))
	}

	src, err := ffile.Open()
	if err != nil {
		log.LogError(err)
		return c.JSON(200, handlers.ResponseError(err))
	}

	defer src.Close()

	// write to mem
	var buf bytes.Buffer

	_, err = io.Copy(&buf, src)
	if err != nil {
		log.LogError(err)
		return c.JSON(200, handlers.ResponseError(err))
	}

	// Check if file is mp3
	dec, _, err := minimp3.DecodeFull(buf.Bytes())
	if err != nil {
		err = errors.New("provided file is not an MP3")
		log.LogError(err)
		return c.JSON(200, handlers.ResponseError(err))
	}

	pretty.Println(dec)

	if dec.Kbps < 12 {
		err = errors.New("provided file is not an MP3")
		log.LogError(err)
		return c.JSON(200, handlers.ResponseError(err))
	}
	if dec.Channels == 0 {
		err = errors.New("provided file is not an MP3")
		log.LogError(err)
		return c.JSON(200, handlers.ResponseError(err))
	}

	if dec.SampleRate == 0 {
		err = errors.New("provided file is not an MP3")
		log.LogError(err)
		return c.JSON(200, handlers.ResponseError(err))
	}

	defer dec.Close()

	// is mp3, save to disk
	err = fileman.SaveFile(ffile.Filename, buf.Bytes())
	if err != nil {
		log.LogError(err)
		return c.JSON(200, handlers.ResponseError(err))
	}

	// Send resposne

	return c.JSON(200, handlers.ResponseError(err))
}

type FileSelect struct {
	Filename string
}

func RemoveFile(c echo.Context) error {
	sel := FileSelect{}
	err := c.Bind(&sel)
	if err != nil {
		log.LogError(err)
		return err
	}

	pretty.Println(sel)

	err = fileman.Remove(sel.Filename)
	c.JSON(200, handlers.ResponseError(err))
	//return fileman.Remove(sel.Filename)
	return err
}
func FileSetDefault(c echo.Context) error {
	sel := FileSelect{}
	err := c.Bind(&sel)
	if err != nil {
		log.LogError(err)
		return err
	}

	// check if this file is valid
	files := fileman.GetFiles()

	var found *fileman.MP3FileInfo

	for _, v := range files {
		v.Selected = false
		if v.Name == sel.Filename {
			found = v
		}
	}
	if found == nil {
		err = errors.New("select file does not exists in file list")
		log.LogError(err)
		return c.JSON(200, handlers.ResponseError(err))
	}

	found.Selected = true

	// set the new default file
	err = playback.PlayFile2(found.AbsolutePath())
	if err != nil {
		log.LogError(err)
		return c.JSON(200, handlers.ResponseError(err))
	}

	// Save settings
	err = settings.SetSelectedFile(found.Name)

	return c.JSON(200, handlers.ResponseError(err))
}

func FileInfo(c echo.Context) error {
	sel := FileSelect{}
	err := c.Bind(&sel)
	if err != nil {
		log.LogError(err)
		return err
	}

	files := fileman.GetFiles()

	pretty.Println(files)

	for _, v := range files {
		if v.Name == sel.Filename {
			return c.JSON(200, v)
		}
	}

	err = errors.New("no such file")
	return c.JSON(200, handlers.ResponseError(err))
}
