package assets

import (
	"github.com/labstack/echo"
)

const _fileIndex = `
<html>
	<head>
		<title>Raspberry MP3LOOP</title>
		<link rel="stylesheet" href="mat.css">
		<script src="mat.js"></script>
		<script src="jquery.js"></script>
		<script src="app.js"></script>
	</head>
	<body>
<div class="container">
   
  <div class="row">
    <h1>RaspberryPI MP3Loop</h1>
  </div>
   <nav class="nav-extended">
      <div class="nav-content">
      <ul class="tabs tabs-transparent">
        <li class="tab"><a class="active" href="#" id="navButtonPlayback">Playback</a></li>
        <li class="tab"><a href="#" id="navButtonFiles">Files</a></li>
        <li class="tab"><a href="#" id="navButtonOutput">OUTPUT device</a></li>
        <li class="tab"><a href="#" id="navButtonLog">Logs</a></li>
      </ul>
    </div>
  </nav>
  <div class="row" id="devices" hidden>
    <h3>Output device</h3>
    <p>This allows you to select the desired output device for the sound to playback from. In case your output device does not show up in the list make sure that the drivers are installed and loaded. Changing the output device will stop the current playback and start it on the newly selected output device.</p>
	<form action="#"><p>
      <label>
        <input type="checkbox" id="outputDeviceForce"  class="filled-in" />
        <span style="color:black;">Force device settings on start.</span>
      </label>
    </p></form>
    <div class="collection" id="outputdevicecollection">
      <a href="#!" class="collection-item">Realtek HD ALC1922 Analog</a>
      <a href="#!" class="collection-item active">Realtek HD ALC1922 Digital-Out</a>
      <a href="#!" class="collection-item">SomeBrand USB DAC</a>
    </div>
    <a class="waves-effect waves-light btn green darken-1" id="outputDeviceScan">Rescan devices</a>


	<div class="row valign-wrapper">
		<div class="col"><h4>Volume</h4></div>
		<div class="col valign-wrapper"><h5><span id="outputDeviceVolVal"></span>%</h5></div>
	</div>
	<div class="row">
	<p>This controls the volume level of the selected output device.</p>
		<p class="range-field">
      <input type="range" id="outputDeviceVolume" min="0" max="100" />
    </p>
	</div>
  </div>
    <div class="row" id="files" hidden>
    <h3>MP3 Files</h3>
    <p>Only one MP3 can be looped at a time. The highlighted file will be looped. Upon changing the desired file, the current playback will be stopped, and the newly selected file will start playing back.</p>
    <div class="collection" id="filescollection">
    </div>
	<h4>Upload new file</h4>
<form id="uploadForm" enctype="multipart/form-data">
<div class="file-field input-field">
      <div class="btn">
        <span>File</span>
        <input type="file" id="fileInput" name="file">
      </div>
      <div class="file-path-wrapper">
        <input class="file-path validate" type="text" placeholder="Upload mp3 files">
      </div>
    </div>
</form>
  </div>
  <div class="row" id="logsDiv" hidden>
	<h3>Logs</h3>
	<p id="logField">Empty</p>
</div>
  <div class="row" id="playback" hidden>
    <h3>Playback</h3>
    <p>This controlls the playback of the current selected mp3 file.</p>
    <div id="playbackinfo">
    <div class="row">
      <a href="#!" class="btn waves-effect waves-red red lighten-2" id="playbackButStop">Stop</a>
      <a class="waves-effect waves-light green lighten-1 btn" id="playbackButRestart">(Re)Start</a>
    </div>
	<div class="row valign-wrapper">
		<div class="col"><h4>Volume</h4></div>
		<div class="col valign-wrapper"><h5><span id="playBackVolVal"></span>%</h5></div>
	</div>
	<div class="row">
	<p>This volume gain is linear and managed by the software, and may greatly reduce the quality of output when you go over 100%. In case the output volume is too low, please use the volume slider under 'Output Device' first, since that slider works on the hardware level and has less change on distortion compared to the slider below.</p>
		<p class="range-field">
      <input type="range" id="playbackVolume" min="0" max="125" />
    </p>
	</div>
        <div class="row">
      <div class="col s4">Status</div>
          <div class="col s8" id="playbackStatus"><span class="green-text">RUNNING</span></div>
    </div>
        <div class="row">
      <div class="col s4">Selected file</div>
          <div class="col s8" id="playbackFile"></div>
    </div>
    <div class="row">
      <div class="col s4">Total lifetime playbacks</div>
      <div class="col s8" id="playBackTotal">17334</div>
    </div>
    </div>
  </div>
</div>
<div id="modalView" class="modal">
    <div class="modal-content">
      <h4>Modal Header</h4>
      <p>A bunch of text</p>
    </div>
    <div class="modal-footer">
      <a href="#!" class="modal-close waves-effect waves-green btn-flat">Agree</a>
    </div>
  </div>
<footer class="page-footer">
  <div class="footer-copyright">
    <div class="container">
      © mrwaggel.be
      <a class="grey-text text-lighten-4 right" href="#!">More Links</a>
    </div>
  </div>
</footer>
	</body>
    <script>
    $( document ).ready(function() {
        AppStart();
    });
    </script>
</html>
`

func HandleIndex(c echo.Context) error {

	return c.HTML(200, _fileIndex)
}
