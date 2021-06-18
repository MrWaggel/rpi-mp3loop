package playcontrols

import (
	"github.com/labstack/echo"
	"mp3loop/easylog"
	"mp3loop/handlers"
	"mp3loop/playback"
	"mp3loop/settings"
	"mp3loop/stats"
)

var log = easylog.NewLogModule("WebPlayControls")

func Initialize(c *echo.Echo) {

	c.GET("/playback/stop", Stop)
	c.GET("/playback/restart", Restart)
	c.GET("/playback/status", Status)
	c.POST("/playback/setvolume", SetVolume)
}

func Stop(c echo.Context) error {
	playback.Stop()

	// Push update of status

	return c.JSON(200, handlers.RespondData(nil))
}

func Restart(c echo.Context) error {
	err := playback.ReStart()

	if err != nil {
		return c.JSON(200, handlers.ResponseError(err))
	}
	// Send playback data
	return c.JSON(200, handlers.RespondData(playback.Status()))
}

type SetVolumePost struct {
	Volume int
}

func SetVolume(c echo.Context) error {
	indata := SetVolumePost{}
	err := c.Bind(&indata)
	if err != nil {
		log.LogError(err)
		return c.JSON(200, handlers.ResponseError(err))
	}

	err = playback.SetVolume(indata.Volume)
	if err != nil {
		log.LogError(err)
		return c.JSON(200, handlers.ResponseError(err))
	}

	// sacve
	err = settings.VolumeSoftwareSet(indata.Volume)
	log.LogError(err)
	return c.JSON(200, handlers.ResponseError(err))

}

func Status(c echo.Context) error {
	// Get volume etc
	vol := playback.VolumeGetInt()

	// status
	isrunning := playback.IsRunning()

	// stats
	statist := stats.Data()

	songname := settings.GetSelectedFile()

	// map it together
	mp := make(map[string]interface{})
	mp["Volume"] = vol
	mp["Stats"] = statist
	mp["Running"] = isrunning
	mp["File"] = songname

	return c.JSON(200, handlers.RespondData(mp))
}
