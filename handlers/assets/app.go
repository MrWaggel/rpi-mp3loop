package assets

import (
	"github.com/labstack/echo"
)

const _file = `function AppStart() {
    Navigation.init();
    Modal.init();

    pagePlayback.init();
    pageFiles.init();
    pageOutput.init();
    pageLog.init();


    pagePlayback.draw();
}

let Navigation = {
    buttons: [],
    init: function () {
        console.log("init nav");
        // Select all buttons

        this.buttons.push($("#navButtonPlayback"));
        this.buttons.push($("#navButtonFiles"));
        this.buttons.push($("#navButtonOutput"));
        this.buttons.push($("#navButtonLog"));

        console.log(this.buttons);
        this.buttons.forEach(function (v, i) {
            console.log(v);
            $(v).on('click', function (evt) {
                Navigation.onClick($(v).attr('id'));
            });
        });
    },
    onClick: function (nameID) {
        console.log("clicked");
        // foreach buttons hide
        Pages.hideAll();

        let selected = null;
        this.buttons.forEach(function (v, i, a) {
            $(v).removeClass("active");
            if (nameID == $(v).attr('id')) {
                selected = v;
            }
        });

        console.log(selected);
        // set active for the nameid
        $(selected).addClass("active");

        let idSelected = selected.attr('id');
        switch (idSelected) {
            case "navButtonPlayback":
                pagePlayback.draw();
                return;
            case "navButtonFiles":
                pageFiles.draw();
                return;
            case "navButtonOutput":
                pageOutput.draw();
                return;
            case "navButtonLog":
                pageLog.draw();
                return;
        }
    }
}

let pagePlayback = {
    selector: null,
    selectors: {
        volume: null,
        volumeSlider: null,
        status: null,
        statsPlays: null,
        file: null,
    },
    lastUpdateData: null,
    init: function () {
        this.selector = $('#playback');

        this.selectors.statsPlays = $('#playBackTotal');
        this.selectors.status = $('#playbackStatus');
        this.selectors.file = $('#playbackFile');
        this.selectors.volumeSlider = $('#playbackVolume');

        Pages.Arr.push(this);
        // hook buttons
        let butStop = $('#playbackButStop');
        butStop.on('click', function (event) {
            // send ajax
            $.get("/playback/stop", function (data) {
                // reload bg
            });
        });
        let butStart = $('#playbackButRestart');
        butStart.on('click', function (event) {
            // send ajax
            $.get("/playback/restart", function (data) {
                // reload bg
            });
        });

        let volVal = $('#playBackVolVal');
        this.selectors.volume = volVal;
        let volumeSlider = $('#playbackVolume');
        volumeSlider.on('input', function () {
            volVal.empty();

            let value = parseInt(volumeSlider.val());

            console.log(value);
            console.log(typeof value);

            volVal.append(value);
            // Update the volume to the server
            $.ajax({
                type: "POST",
                url: "/playback/setvolume",
                async: false,
                data: JSON.stringify({Volume: value}),
                contentType: "application/json",
                complete: function (data) {
                    console.log("ajax done");
                    let resp = data.responseJSON;
                }
            });
        });

        this.backgroundUpdate();
    },
    draw: function () {
        Pages.hideAll();
        this.selector.show();
    },
    backgroundUpdate: function () {
        // Poll and recall
        $.get("/playback/status", function (pdata) {
            if (pdata.Good) {
                let data = pdata.Data;
                let old = pagePlayback.lastUpdateData;
                // update
                if (this.lastUpdateData == null) {
                    // Just write
                    pagePlayback.selectors.volume.html(data.Volume);
                    pagePlayback.selectors.volumeSlider.val(data.Volume);
                    pagePlayback.updateStatus(data.Running);
                    pagePlayback.selectors.file.html(data.File);
                    pagePlayback.selectors.statsPlays.html(data.Stats.StartTotalPlays);
                } else {
                    if (data.Volume != old.Volume) {
                        pagePlayback.selectors.volume.html(data.Volume);
                        pagePlayback.selectors.volumeSlider.val(data.Volume);

                    }
                    if (data.Runnning != old.Running) {
                        pagePlayback.updateStatus(data.Runnning);
                    }
                    if (data.File != old.File) {
                        pagePlayback.selectors.file.html(data.File);
                    }
                    if (pdata.Stats.StartTotalPlays != old.Stats.StartTotalPlays) {
                        pagePlayback.selectors.statsPlays.html(data.Stats.StartTotalPlays);
                    }
                }

                pagePlayback.lastUpdateData = pdata;
            }
        }).always(function () {
            setTimeout(function () {
                pagePlayback.backgroundUpdate();
                // wait a minute after you recieved the data
            }, 1000);
        });
    },
    updateStatus: function (bool) {
        this.selectors.status.empty();
        let status = "STOPPED";
        let color = "red"
        if (bool) {
            status = "RUNNING"
            color = "green";
        }

        let sp = $('<span>')
        sp.addClass(color + "-text");
        sp.append(status)
        this.selectors.status.append(sp);
    }
}

let pageFiles = {
    selector: null,
    collection: null,
    init: function () {
        this.selector = $('#files');
        this.collection = $('#filescollection');
        Pages.Arr.push(this);

        // File upload handling
        $('#fileInput').on('change', function () {
            var file = this.files[0];
            pageFiles.uploadFile(file);
        });
    },
    uploadFile: function (file) {
        // Create modal
        let m = easyModal();
        m.setTitle("Uploading file");

        // Create body stuff
        let b = $('<div>')
        let p = $('<p>Uploading file, please don\'t close the dialog.</p>')
        let progr = $('<progress></progress>');
        b.append(p);
        b.append(progr);
        m.setBody(b);

        m.setOptions({dismissible: false});

        // pop it
        let modalControl = m.draw();

        console.log($('#uploadForm'));

        // Do ajax shit
        $.ajax({
            // Your server script to process the upload
            url: '/files/add',
            type: 'POST',

            // Form data
            data: new FormData($('#uploadForm')[0]),

            // Tell jQuery not to process data or worry about content-type
            // You *must* include these options!
            cache: false,
            contentType: false,
            processData: false,
            success: function (rda) {
                modalControl.forEach(function (elem) {
                    elem.close();
                });
                // Create new modal
                console.log(rda);
                let resp = rda;
                if (resp.Good == true) {
                    // Show succes modal and reload background
                    let m = easyModal();
                    m.setTitle("File uploaded!");
                    let d = $('<div>');
                    d.append($('<p>The file has been successfully uploaded.</p>'))
                    m.setBody(d);
                    m.draw();

                    // reload files
                    pageFiles.draw();

                } else if (resp.Good == false) {
                    // Show error
                    let m = easyModal();
                    m.setTitle("File upload failed");
                    let d = $('<div>');
                    d.append($('<p>The upload of the file failed. The following reason was given</p>'))
                    d.append($('<p>' + resp.Error + '</p>'))
                    m.setBody(d);
                    m.draw();
                }

            },

            // Custom XMLHttpRequest
            xhr: function () {
                var myXhr = $.ajaxSettings.xhr();
                if (myXhr.upload) {
                    // For handling the progress of the upload
                    myXhr.upload.addEventListener('progress', function (e) {
                        console.log(e);
                        if (e.lengthComputable) {
                            progr.attr({
                                value: e.loaded,
                                max: e.total,
                            });
                        }
                    }, false);
                }
                return myXhr;
            }
        });

    },
    draw: function () {
        Pages.hideAll();

        this.collection.empty();
        this.selector.show();

        // Get all files
        $.get("files/all", function (data) {
            data.forEach(function (v) {
                pageFiles.collectionAdd(v);
            });
        });

    },
    collectionAdd: function (data) {
        let a = $('<a href="#!" class="collection-item">');
        a.on('click', function (event) {
            // create popup modal
            let m = easyModal();
            m.setTitle("Change playback file");
            let p = $('<p>')
            p.append("Are you sure you want to change the playback to <strong>" + data.Name + "</strong>?");

            // button
            let delbut = $('<a href="#!" class="modal-close waves-effect waves-green btn">Change</a>');
            delbut.on('click', function (event) {

                $.ajax({
                    type: "POST",
                    url: "/files/use",
                    async: false,
                    data: JSON.stringify({Filename: data.Name}),
                    contentType: "application/json",
                    complete: function (data) {
                        console.log("ajax done");
                        let resp = data.responseJSON;
                        if (resp.Good == true) {
                            pageFiles.draw();
                        } else {
                            // fail
                        }
                    }
                });
            });
            m.setBody(p);
            m.addButton(delbut);
            m.draw();
        });

        if (data.Selected == true) {
            a.addClass("active");
        }
        // add the name
        a.append(data.Name);

        let secDiv = $('<div class="secondary-content">');
        // add info and delte button
        let info = $('<span>info</span>');
        info.on('click', function (event) {
            event.stopPropagation();
            // Get file info
            console.log(data);


            $.ajax({
                type: "POST",
                url: "files/info",
                async: false,
                data: JSON.stringify({Filename: data.Name}),
                contentType: "application/json",
                complete: function (data) {
                    console.log("ajax done");
                    console.log(data.responseJSON);
                    let m = easyModal();
                    m.setTitle("File info");
                    m.setBody(easyTable(data.responseJSON))
                    m.draw();
                }
            });
        });
        let del = $('<span>delete</span>');
        del.on('click', function (event) {
            // Get file info
            event.stopPropagation();

            let m = easyModal();
            m.setTitle("Delete file");
            let p = $('<p>')
            p.append("Are you sure you want to remove the file named <strong>" + data.Name + "</strong>?");

            // button
            let delbut = $('<a href="#!" class="modal-close waves-effect waves-green btn">Delete</a>');
            delbut.on('click', function (event) {

                $.ajax({
                    type: "POST",
                    url: "files/delete",
                    async: false,
                    data: JSON.stringify({Filename: data.Name}),
                    contentType: "application/json",
                    complete: function (data) {
                        console.log("ajax done");
                        let resp = (data.responseJSON);
                        if (resp.Good == true) {
                            pageFiles.draw();
                        }
                    }
                });
            });
            m.setBody(p);
            m.addButton(delbut);
            m.draw();

        });

        secDiv.append(info);
        secDiv.append(" | ");
        secDiv.append(del);
        a.append(secDiv);

        this.collection.append(a);
    },
}

let easyTable = function (objec) {
    let t = $('<table class="highlight">');

    let b = $('<tbody>');
    t.append(b);

    console.log(objec);

    $.each(objec, function (name, val) {
        let tr = $('<tr>');
        tr.append($('<td>' + name + '</td>'));
        tr.append($('<td>' + val + '</td>'));
        b.append(tr);
    });
    return t;
}

let easyModal = function () {
    out = {
        title: null,
        options: null,
        body: null,
        buttons: [],
        addButton: function (button) {
            this.buttons.push(button);
        },
        setBody: function (html) {
            this.body = html;
        },
        setTitle: function (title) {
            this.title = title;
        },
        setOptions: function (opts) {
            this.options = opts;
        },
        draw: function () {
            let wrap = $('<div>');
            let content = $('<div class="modal-content">');
            content.append($('<h4>' + this.title + '</h4>'));

            content.append(this.body);

            wrap.append(content);

            if (this.buttons.length > 0) {
                let bcont = $('<div class="modal-footer">');
                this.buttons.forEach(function (v) {
                    bcont.append(v);
                });
                wrap.append(bcont);

            }

            return Modal.draw(wrap, this.buttons);
        }
    }

    return out;
}

let Modal = {
    selector: null,
    init: function () {
        this.selector = $('#modalView');
    },
    draw: function (innerhtml, options) {
        this.selector.empty();
        this.selector.append(innerhtml);
        let instance = M.Modal.init(this.selector, options);
        instance.forEach(function (v) {
            v.open();
        });

        return instance;
    }
}

let pageLog = {
    selector: null,
    selectors: {
        div: null,
        paragraph: null,
    },
    init: function () {
        this.selectors.paragraph = $('#logField');
        this.selectors.div = $('#logsDiv');
        this.selector = this.selectors.div;
        Pages.Arr.push(this);
    },
    draw: function () {
        // do get call
        this.selectors.div.show();
        $.get("logs", function (fd) {
            // get data array
            let darr = fd.Data;
            pageLog.selectors.paragraph.empty();
            darr.forEach(function (val) {
                pageLog.writeLine(val);
            });
        });
    },
    writeLine: function (line) {
        this.selectors.paragraph.append(line + "<br>")
    },
}

let pageOutput = {
    selector: null,
    collection: null,
    selectors: {
        forceCheckBox: null,
        volumeSlider: null,
        volumeValue: null,
        scanButton: null,
    },
    init: function () {
        this.selector = $('#devices');
        this.collection = $('#outputdevicecollection');

        this.selectors.forceCheckBox = $('#outputDeviceForce');
        this.selectors.forceCheckBox.change(function () {
            // set get request
            $.get("devices/settings/force/toggle", function (fd) {
                let forceCheckBox = pageOutput.selectors.forceCheckBox;
                let val = forceCheckBox.prop("checked");
                if (val == null) {
                    val = false;
                }
                // get and swap  current val
                //forceCheckBox.prop( "checked", !val);
            });
        });
        this.selectors.volumeSlider = $('#outputDeviceVolume');
        this.selectors.volumeSlider.on('input', function () {
            let volumeSlider = pageOutput.selectors.volumeSlider;
            let volValue = pageOutput.selectors.volumeValue;
            volValue.empty();

            let value = parseInt(volumeSlider.val());

            volValue.append(value);
            // Update the volume to the server
            $.ajax({
                type: "POST",
                url: "/devices/settings/volume",
                async: false,
                data: JSON.stringify({Volume: value}),
                contentType: "application/json",
                complete: function (data) {
                }
            });
        });
        this.selectors.volumeValue = $('#outputDeviceVolVal');
        this.selectors.scanButton = $('#outputDeviceScan');
        this.selectors.scanButton.on('click', function (event) {
            pageOutput.draw();
        });
        Pages.Arr.push(this);
    },
    draw: function () {
        // Get all files
        $.get("devices/all", function (fd) {
            let devs = fd["Devices"];
            devs.forEach(function (v) {
                pageOutput.collectionAdd(v);
            });

            // set volume and checkbox
            let po = pageOutput;
            po.selectors.forceCheckBox.prop("checked", fd["Force"]);
            po.selectors.volumeSlider.val(fd["Volume"]);
            po.selectors.volumeValue.html(fd["Volume"]);
        });

        // Get settings

        this.collection.empty();
        Pages.hideAll();
        this.selector.show();
    },
    collectionAdd(data) {
        let a = $('<a href="#!" class="collection-item">');


        if (data.Selected == true) {
            a.addClass("active");
        }
        // add the name
        a.append("ID: " + data.IDSink + "")
        a.append(data.NameAlsa);

        let secDiv = $('<div class="secondary-content">');
        // add info and delte button
        let info = $('<span>info</span>');
        info.on('click', function (event) {
            event.stopPropagation();
            // Get file info
            console.log("clicked aa");
            console.log("clicked aa");
            console.log("clicked aa");
            $.ajax({
                type: "POST",
                url: "devices/info",
                async: false,
                data: JSON.stringify({DeviceID: data.IDSink}),
                contentType: "application/json",
                complete: function (data) {
                    console.log("ajax info");
                    console.log(data.responseJSON);
                    let modal3 = easyModal();
                    modal3.setTitle("Device info info");
                    modal3.setBody(easyTable(data.responseJSON.Data));
                    let mdraw = modal3.draw();
                }
            });
        });


        a.on('click', function (event) {
            // Popup
            let modal = easyModal();
            modal.setTitle("Change output device.");
            modal.setBody($('<p>Are you sure you want to change the audio device?</p>'));

            // Buttons
            let but = $('<a href="#!" class="modal-close waves-effect waves-green btn">Change</a>');
            modal.addButton(but);
            let modalcontrol = modal.draw();
            but.on('click', function (event) {
                console.log("clicked useeee");
                console.log("clicked useeee");
                console.log("clicked SENDING AJX");
                // send jqeruy
                $.ajax({
                    type: "POST",
                    url: "/devices/use",
                    async: false,
                    data: JSON.stringify({DeviceID: data.IDSink}),
                    contentType: "application/json",
                    complete: function (data) {
                        modalcontrol.close();
                        console.log("ajax done");
                        console.log("ajax done");
                        let resp = (data.responseJSON);
                        if (resp.Good == true) {
                            pageFiles.draw();
                        } else {
                            // show fail
                            let m = easyModal();
                            m.setTitle("Device change failed");
                            m.setBody($('<p>We failed to change the device, reason: ' + resp.Error + '</p>'))
                        }
                    }
                });
            });
        });

        secDiv.append(info);
        a.append(secDiv);

        this.collection.append(a);
    },
}


let Pages = {
    Arr: [],
    hideAll: function () {
        this.Arr.forEach(function (v) {
            v.selector.hide();
        });
    },
}`

func HandleAppJS(c echo.Context) error {
	return c.String(200, _file)
}

func HandleAppJSStatic(c echo.Context) error {
	return c.File("/home/dev/go/src/mp3loop/handlers/assets/app.ja.js")
}
