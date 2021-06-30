package main

import (
	"github.com/labstack/echo"
	"mp3loop/devices"
	"mp3loop/easylog"
	"mp3loop/fileman"
	"mp3loop/handlers"
	"mp3loop/handlers/assets"
	"mp3loop/handlers/devicecontrol"
	"mp3loop/handlers/filecontrols"
	"mp3loop/handlers/playcontrols"
	"mp3loop/playback"
	"mp3loop/settings"
	"mp3loop/stats"
)

func main() {

	settings.Initialize()
	stats.Initialize()
	devices.Initialize()
	fileman.Initialize()
	playback.Initialize()

	e := echo.New()

	filecontrols.Initialize(e)
	devicecontrol.Initialize(e)
	playcontrols.Initialize(e)

	// Static files
	e.GET("/", assets.HandleIndex)
	e.GET("/mat.css", assets.HandlePureCSS)
	e.GET("/mat.js", assets.HandleMaterializeJS)
	e.GET("/jquery.js", assets.HandleJqueryJS)
	e.GET("/app.js", assets.HandleAppJS)
	e.GET("/logs", func(context echo.Context) error {
		// get log
		l := easylog.GetLog()
		return context.JSON(200, handlers.RespondData(l))
	})

	e.Start(":8080")
}
