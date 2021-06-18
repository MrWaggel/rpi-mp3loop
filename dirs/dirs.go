package dirs

import (
	"github.com/kardianos/osext"
	"os"
	"path"
)

var _execFolder string

func init() {
	// Create data dir
	_execFolder, _ = osext.ExecutableFolder()
	makeFolder(DirData())
	makeFolder(DirFiles())
}

func makeFolder(f string) error {
	if _, err := os.Stat("/path/to/whatever"); os.IsNotExist(err) {
		// path/to/whatever does *not* exist
		return os.Mkdir(f, os.ModePerm)
	}
	return nil
}

func DirData() string {
	return path.Join(_execFolder, "data")
}

func DirFiles() string {
	return path.Join(DirData(), "files")
}
