package devicecontrol

import (
	"errors"
	"github.com/labstack/echo"
	"mp3loop/devices"
	"mp3loop/easylog"
	"mp3loop/handlers"
	"mp3loop/settings"
)

func Initialize(c *echo.Echo) {
	c.GET("/devices/all", GetDevices)
	c.POST("/devices/info", DeviceInfo)
	c.POST("/devices/use", SetDefaultDevice)
	c.GET("/devices/settings/force/toggle", SettingsForceToggle)
	c.POST("/devices/settings/volume", SettingsVolume)
}

var log = easylog.NewLogModule("WebDeviceControl")

func SettingsForceToggle(c echo.Context) error {
	// get current value
	v := settings.ForceOutputDeviceOnStartGet()
	v = !v

	err := settings.ForceOutputDeviceOnstartSet(v)
	if err != nil {
		log.LogError(err)
		return c.JSON(200, handlers.ResponseError(err))
	}

	// Send new value
	return c.JSON(200, handlers.RespondData(v))
}

func GetDevices(c echo.Context) error {
	// get list devices
	devs, err := devices.GetDevices()
	if err != nil {
		return c.JSON(200, handlers.ResponseError(err))
	}

	// get volume of selected device
	vol := 0
	for _, v := range devs {
		if v.Selected {
			vol = v.GetVolumePercentage()
			break
		}
	}
	// Get force checkbox
	force := settings.ForceOutputDeviceOnStartGet()

	mp := make(map[string]interface{})
	mp["Devices"] = devs
	mp["Force"] = force
	mp["Volume"] = vol
	return c.JSON(200, mp)
}

type RequestVolume struct {
	Volume int
}

func SettingsVolume(c echo.Context) error {
	// Get new val
	vs := RequestVolume{}
	err := c.Bind(&vs)
	if err != nil {
		log.LogError(err)
		return c.JSON(200, handlers.ResponseError(err))
	}
	// Change in devices
	err = devices.SetVolume(vs.Volume)
	if err != nil {
		return c.JSON(200, handlers.ResponseError(err))
	}
	// save to settings

	err = settings.VolumeDeviceSet(vs.Volume)
	if err != nil {
		return c.JSON(200, handlers.ResponseError(err))
	}
	return c.JSON(200, handlers.RespondData(vs.Volume))
}

type RequestDevice struct {
	DeviceID int
}

func DeviceInfo(c echo.Context) error {
	rd := RequestDevice{}
	err := c.Bind(&rd)
	if err != nil {
		log.LogError(err)
		return c.JSON(200, handlers.ResponseError(err))
	}

	// Check if that device exists
	devs, err := devices.GetDevices()
	if err != nil {
		log.LogError(err)
		return c.JSON(200, handlers.ResponseError(err))
	}

	for _, v := range devs {
		if v.IDSink == rd.DeviceID {
			// send this info
			return c.JSON(200, handlers.RespondData(v))
		}
	}

	err = errors.New("device id doesn't exist")
	log.LogError(err)
	return c.JSON(200, handlers.ResponseError(err))

}
func SetDefaultDevice(c echo.Context) error {
	rd := RequestDevice{}
	err := c.Bind(&rd)
	if err != nil {
		log.LogError(err)
		return c.JSON(200, handlers.ResponseError(err))
	}

	// Check if that device exists
	devs, err := devices.GetDevices()
	if err != nil {
		log.LogError(err)
		return c.JSON(200, handlers.ResponseError(err))
	}

	devname := ""
	for _, v := range devs {
		if v.IDSink == rd.DeviceID {
			devname = v.NameDevice
			break
		}
	}

	if devname == "" {
		err = errors.New("device id doesn't exist")
		log.LogError(err)
		return c.JSON(200, handlers.ResponseError(err))
	}

	// swap
	err = devices.SetDefaultDevice(devname)
	if err != nil {
		log.LogError(err)
		return c.JSON(200, handlers.ResponseError(err))
	}

	// Save to settings
	err = settings.DeviceOutputDefaultSet(devname)

	return c.JSON(200, handlers.ResponseError(err))
}
